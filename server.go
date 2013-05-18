package main

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "log"
    "os/exec"
)

func WriteGitToHttp(w http.ResponseWriter, args ...string) {
    cmd := exec.Command("git", args...)
    stdout, err := cmd.StdoutPipe()

    if err != nil {
        w.WriteHeader(404)
        log.Fatal("Error:", err)
        return
    }

    if err := cmd.Start(); err != nil {
        w.WriteHeader(404)
        log.Fatal("Error:", err)
        return
    }

    var n int
    for i := 0; i < 32; i++ {
        data := make([]byte, 256)
        nbytes, e := stdout.Read(data)

        if nbytes == 0 || e != nil {
            break
        }
        n += nbytes
        w.Write(data)
    }
    log.Printf("Sent %d bytes to client.", n)
}

func getInfoRefs(route *Route, w http.ResponseWriter, r *http.Request) {
    log.Printf("getInfoRefs for %s", route.RepoPath)
    // TODO: find repo path at route.RepoPath
    WriteGitToHttp(w, "upload-pack", "--stateless-rpc", "--advertise-refs", ".")
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
    fmt.Printf("Hi there, I love %s!\n", r.URL.Path[0:])

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

