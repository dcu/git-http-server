package gitserver

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// RouteFunc defines the prototype of a route handler function.
type RouteFunc func(route *Route, w http.ResponseWriter, r *http.Request)

// RouteMatcher has a regexp to match the route and a handler for that route.
type RouteMatcher struct {
	Matcher *regexp.Regexp
	Handler RouteFunc
}

// Route has the repository file with the matched route.
type Route struct {
	RepoPath     string
	File         string
	MatchedRoute RouteMatcher
}

// Dispatch processes the incoming http request.
func (route *Route) Dispatch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	route.MatchedRoute.Handler(route, w, r)
}

// NewParsedRoute returns a new instance of a Route
func NewParsedRoute(repoName string, file string, matcher RouteMatcher) *Route {
	return &Route{RepoPath: repoName, File: file, MatchedRoute: matcher}
}

// Routes contains a list of the known routes to be handled.
var Routes = []RouteMatcher{
	RouteMatcher{Matcher: regexp.MustCompile("(.*?)/info/refs$"), Handler: getInfoRefs},
	RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-upload-pack$"), Handler: uploadPack},
	RouteMatcher{Matcher: regexp.MustCompile("(.*?)/git-receive-pack$"), Handler: receivePack},
}

// MatchRoute returns the matched route or nil.
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
