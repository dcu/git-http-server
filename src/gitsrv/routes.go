package gitsrv

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "log"
)

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
    r.ParseForm()
    route.MatchedRoute.Handler(route, w, r)
}

func NewParsedRoute(repoName string, file string, matcher RouteMatcher) *Route {
    return &Route{RepoPath: repoName, File: file, MatchedRoute: matcher};
}

var Routes = []RouteMatcher{
    RouteMatcher{Matcher: regexp.MustCompile("(.*?)/info/refs$"), Handler: getInfoRefs},
    RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-upload-pack$"), Handler: uploadPack},
    RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-receive-pack$"), Handler: receivePack},
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

