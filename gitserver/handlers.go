package gitserver

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dcu/go-authy"
)

var (
	oneTouchExpiresAfter = 45 * time.Second
)

// AbsoluteRepoPath returns the absolute path for the given relative repository path
func AbsoluteRepoPath(relativePath string) (string, error) {
	if !strings.HasSuffix(relativePath, ".git") {
		relativePath += ".git"
	}

	path := fmt.Sprintf("%s/%s", gServerConfig.Repos.Path, relativePath)
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if strings.Contains(path, "..") {
		return "", errors.New("invalid repo path")
	}

	return absolutePath, nil
}

func getInfoRefs(route *Route, w http.ResponseWriter, r *http.Request) {
	repo, err := AbsoluteRepoPath(route.RepoPath)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if gServerConfig.Repos.AutoInit && !repoExists(repo) {
		cmd := GitCommand{Args: []string{"init", "--bare", repo}}
		_, err := cmd.Run(true)
		if err != nil {
			w.WriteHeader(404)
			return
		}
	}

	log.Printf("getInfoRefs for %s", repo)

	serviceName := getServiceName(r)

	message := messageFromService(serviceName, route.RepoPath)
	details := authy.Details{
		"repo": repo,
		"ip":   r.RemoteAddr,
	}
	if !approveTransaction(message, details) {
		w.WriteHeader(403)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/x-git-"+serviceName+"-advertisement")

	str := "# service=git-" + serviceName
	fmt.Fprintf(w, "%.4x%s\n", len(str)+5, str)
	fmt.Fprintf(w, "0000")
	WriteGitToHTTP(w, GitCommand{Args: []string{serviceName, "--stateless-rpc", "--advertise-refs", repo}})
}

func getServiceName(r *http.Request) string {
	if len(r.Form["service"]) > 0 {
		return strings.Replace(r.Form["service"][0], "git-", "", 1)
	}

	return ""
}

func uploadPack(route *Route, w http.ResponseWriter, r *http.Request) {
	repo, err := AbsoluteRepoPath(route.RepoPath)
	if err != nil {
		return
	}
	log.Printf("uploadPack for %s", repo)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/x-git-upload-pack-result")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(404)
		log.Fatal("Error:", err)
		return
	}

	WriteGitToHTTP(w, GitCommand{ProcInput: bytes.NewReader(requestBody), Args: []string{"upload-pack", "--stateless-rpc", repo}})
}

func receivePack(route *Route, w http.ResponseWriter, r *http.Request) {
	repo, err := AbsoluteRepoPath(route.RepoPath)
	if err != nil {
		return
	}
	log.Printf("receivePack for %s", repo)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/x-git-receive-pack-result")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(404)
		log.Fatal("Error:", err)
		return
	}

	WriteGitToHTTP(w, GitCommand{ProcInput: bytes.NewReader(requestBody), Args: []string{"receive-pack", "--stateless-rpc", repo}})
}

func goGettable(route *Route, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")

	url := fmt.Sprintf("%s%s", r.Host, r.URL.Path)
	fmt.Fprintf(w, `<html><head><meta name="go-import" content="%s git https://%s"></head><body>go get %s</body></html>`, url, url, url)
}

func repoExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func approveTransaction(message string, details authy.Details) bool {
	if !gServerConfig.HasAuthy() {
		// Ignore request.
		return true
	}

	authyAPI := authy.NewAuthyAPI(gServerConfig.Authy.APIKey)
	request, err := authyAPI.SendApprovalRequest(
		gServerConfig.Authy.UserID,
		message,
		details,
		url.Values{
			"seconds_to_expire": {strconv.FormatInt(int64(oneTouchExpiresAfter), 10)},
		},
	)
	if err != nil {
		return false
	}

	status, err := authyAPI.WaitForApprovalRequest(request.UUID, oneTouchExpiresAfter, url.Values{})
	if err != nil {
		return false
	}

	return status == authy.OneTouchStatusApproved
}

func messageFromService(service string, repo string) string {
	message := ""
	if service == "receive-pack" {
		message = "Push to " + repo
	} else if service == "upload-pack" {
		message = "Fetch from " + repo
	} else {
		message = "Unknown service " + service + " for " + repo
	}

	return message
}
