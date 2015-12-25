package gitserver

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// AbsoluteRepoPath returns the absolute path for the given relative repository path
func AbsoluteRepoPath(relativePath string) (string, error) {
	if !strings.HasSuffix(relativePath, ".git") {
		relativePath += ".git"
	}

	path := fmt.Sprintf("%s/%s", gServerConfig.ReposRoot, relativePath)
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if strings.Contains(path, "..") {
		return "", errors.New("Invalid repo path.")
	}

	return absolutePath, nil
}

func getInfoRefs(route *Route, w http.ResponseWriter, r *http.Request) {
	repo, err := AbsoluteRepoPath(route.RepoPath)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if gServerConfig.AutoInitRepos && !repoExists(repo) {
		cmd := GitCommand{Args: []string{"init", "--bare", repo}}
		_, err := cmd.Run(true)
		if err != nil {
			w.WriteHeader(404)
			return
		}
	}

	log.Printf("getInfoRefs for %s", repo)

	serviceName := getServiceName(r)
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
