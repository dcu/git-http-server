package gitserver

import (
	"net/http"
)

// Middleware tries to match the route with the git server otherwise it calls to "next"
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedRoute := MatchRoute(r)
		if parsedRoute != nil {
			parsedRoute.Dispatch(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// MiddlewareFunc tries to match the route with the git server otherwise it calls to "next"
func MiddlewareFunc(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parsedRoute := MatchRoute(r)
		if parsedRoute != nil {
			parsedRoute.Dispatch(w, r)
		} else {
			next(w, r)
		}
	}
}
