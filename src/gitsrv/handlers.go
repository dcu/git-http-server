package gitsrv

import (
    "fmt"
    "net/http"
    "strings"
    "log"
    "io/ioutil"
    "bytes"
)

func getInfoRefs(route *Route, w http.ResponseWriter, r *http.Request) {
    log.Printf("getInfoRefs for %s", route.RepoPath)
    // TODO: find repo path at route.RepoPath

    serviceName := getServiceName(r)
    w.Header().Set("Content-Type", "application/x-git-"+serviceName+"-advertisement")

    str := "# service=git-"+serviceName
    fmt.Fprintf(w, "%.4x%s\n", len(str)+5, str )
    fmt.Fprintf(w, "0000")
    WriteGitToHttp(w, GitCommand{args: []string{serviceName, "--stateless-rpc", "--advertise-refs", route.RepoPath}} )
}

func getServiceName(r *http.Request) string {
    if len(r.Form["service"]) > 0 {
        return strings.Replace(r.Form["service"][0], "git-", "", 1)
    }

    return ""
}

func uploadPack(route *Route, w http.ResponseWriter, r *http.Request) {
    log.Printf("uploadPack for %s", route.RepoPath)
    // TODO: find repo path at route.RepoPath

    w.Header().Set("Content-Type", "application/x-git-upload-pack-result")

    requestBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(404)
        log.Fatal("Error:", err)
        return
    }

    WriteGitToHttp(w, GitCommand{procInput: bytes.NewReader(requestBody), args: []string{"upload-pack", "--stateless-rpc", route.RepoPath}})
}

func receivePack(route *Route, w http.ResponseWriter, r *http.Request) {
    log.Printf("receivePack for %s", route.RepoPath)
    // TODO: find repo path at route.RepoPath

    w.Header().Set("Content-Type", "application/x-git-receive-pack-result")

    requestBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(404)
        log.Fatal("Error:", err)
        return
    }

    WriteGitToHttp(w, GitCommand{procInput: bytes.NewReader(requestBody), args: []string{"receive-pack", "--stateless-rpc", route.RepoPath}})
}


