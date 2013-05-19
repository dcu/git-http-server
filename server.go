package main

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "log"
    "os/exec"
    "io"
    "io/ioutil"
    "bytes"
)

type GitCommand struct {
    procInput *bytes.Reader
    args []string
}

func WriteGitToHttp(w http.ResponseWriter, gitCommand GitCommand) {
    cmd := exec.Command("git", gitCommand.args...)
    stdout, err := cmd.StdoutPipe()

    if err != nil {
        w.WriteHeader(404)
        log.Fatal("Error:", err)
        return
    }

    if gitCommand.procInput != nil {
        cmd.Stdin = gitCommand.procInput
    }


    if err := cmd.Start(); err != nil {
        w.WriteHeader(404)
        log.Fatal("Error:", err)
        return
    }

    nbytes, err := io.Copy(w, stdout)
    log.Printf("Bytes written: %d", nbytes)
    if err != nil {
        log.Fatal("Error writing to socket", err)
    }
}


func getInfoRefs(route *Route, w http.ResponseWriter, r *http.Request) {
    log.Printf("getInfoRefs for %s", route.RepoPath)
    // TODO: find repo path at route.RepoPath

    w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")

    str := "# service=git-upload-pack"
    fmt.Fprintf(w, "%.4x%s\n", len(str)+5, str )
    fmt.Fprintf(w, "0000")
    WriteGitToHttp(w, GitCommand{args: []string{"upload-pack", "--stateless-rpc", "--advertise-refs", "."}} )
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

    WriteGitToHttp(w, GitCommand{procInput: bytes.NewReader(requestBody), args: []string{"upload-pack", "--stateless-rpc", "."}})
}


func getRepoFile(route *Route, w http.ResponseWriter, r *http.Request) {
    log.Printf("getRepoFile for %s", route.RepoPath)

    http.ServeFile(w, r, ".git/"+route.File)
}


type RouteFunc func (route *Route, w http.ResponseWriter, r *http.Request)
type RouteMatcher struct {
    Matcher *regexp.Regexp
    Handler RouteFunc
}

type Route struct {
    RepoPath string
    File string
    MatchedRoute RouteMatcher
}

func (route *Route) Dispatch(w http.ResponseWriter, r *http.Request) {
    route.MatchedRoute.Handler(route, w, r)
}

func NewParsedRoute(repoName string, file string, matcher RouteMatcher) *Route {
    return &Route{RepoPath: repoName, File: file, MatchedRoute: matcher};
}

var Routes = []RouteMatcher{
    RouteMatcher{Matcher: regexp.MustCompile("(.*?)/info/refs$"), Handler: getInfoRefs},
    RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-upload-pack$"), Handler: uploadPack},
    RouteMatcher{Matcher: regexp.MustCompile("(.*?)/HEAD$"), Handler: getRepoFile},
}

func MatchRoute(r *http.Request) *Route {
    path := r.URL.Path[1:]

    for _, routeHandler := range Routes {
        matches := routeHandler.Matcher.FindStringSubmatch(path)
        if matches != nil {
            repoName := matches[1]
            file := strings.Replace(path, repoName+"/", "", 1)

            fmt.Printf("matches: %q\n", matches)
            fmt.Printf("repo name: %s\n", repoName)
            fmt.Printf("file: %s\n", file)

            return NewParsedRoute(repoName, file, routeHandler)
        }
    }

    log.Printf("No route found for: %s", path)
    return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Processing %s!\n", r.URL.Path[0:])

    parsedRoute := MatchRoute(r)
    if parsedRoute != nil {
        parsedRoute.Dispatch(w, r)
    }
}

func main() {
    log.Printf("Starting server on localhost:4000")
    http.HandleFunc("/", handler)
    http.ListenAndServe(":4000", nil)
}

