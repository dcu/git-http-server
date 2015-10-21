package main

import (
	"flag"
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"log"
	"net/http"
)

var (
	listenAddressFlag = flag.String("web.listen-address", ":4000", "Address on which to listen to git requests.")
	authUserFlag      = flag.String("auth.user", "", "Username for basic auth.")
	authPassFlag      = flag.String("auth.pass", "", "Password for basic auth.")
)

func authMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || password != *authPassFlag || user != *authUserFlag {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"metrics\"")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			next(w, r)
		}
	}
}

func hasUserAndPassword() bool {
	return *authUserFlag != "" && *authPassFlag != ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "nothing to see here\n")
}

func main() {
	flag.Parse()

	log.Printf("Starting server on localhost:4000")
	app := gitserver.MiddlewareFunc(handler)
	if hasUserAndPassword() {
		app = authMiddleware(app)
	}
	http.HandleFunc("/", app)
	http.ListenAndServe(*listenAddressFlag, nil)
}
