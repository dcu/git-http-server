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

func absoluteRepoPath(relativePath string) (string, error) {
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
	repo, err := absoluteRepoPath(route.RepoPath)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if gServerConfig.AutoInitRepos && !repoExists(repo) {
		cmd := GitCommand{args: []string{"init", "--bare", repo}}
		_, err := cmd.Run(true)
		if err != nil {
			w.WriteHeader(404)
			return
		}
	}

	log.Printf("getInfoRefs for %s", repo)

	serviceName := getServiceName(r)
	w.Header().Set("Content-Type", "application/x-git-"+serviceName+"-advertisement")

	str := "# service=git-" + serviceName
	fmt.Fprintf(w, "%.4x%s\n", len(str)+5, str)
	fmt.Fprintf(w, "0000")
	WriteGitToHTTP(w, GitCommand{args: []string{serviceName, "--stateless-rpc", "--advertise-refs", repo}})
}

func getServiceName(r *http.Request) string {
	if len(r.Form["service"]) > 0 {
		return strings.Replace(r.Form["service"][0], "git-", "", 1)
	}

	return ""
}

func uploadPack(route *Route, w http.ResponseWriter, r *http.Request) {
	repo, err := absoluteRepoPath(route.RepoPath)
	if err != nil {
		return
	}
	log.Printf("uploadPack for %s", repo)

	w.Header().Set("Content-Type", "application/x-git-upload-pack-result")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(404)
		log.Fatal("Error:", err)
		return
	}

	WriteGitToHTTP(w, GitCommand{procInput: bytes.NewReader(requestBody), args: []string{"upload-pack", "--stateless-rpc", repo}})
}

func receivePack(route *Route, w http.ResponseWriter, r *http.Request) {
	repo, err := absoluteRepoPath(route.RepoPath)
	if err != nil {
		return
	}
	log.Printf("receivePack for %s", repo)

	w.Header().Set("Content-Type", "application/x-git-receive-pack-result")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(404)
		log.Fatal("Error:", err)
		return
	}

	WriteGitToHTTP(w, GitCommand{procInput: bytes.NewReader(requestBody), args: []string{"receive-pack", "--stateless-rpc", repo}})
}

func repoExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
