package main

import (
	"flag"
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"log"
	"net/http"
	"os"
)

var (
	listenAddressFlag = flag.String("web.listen-address", ":4000", "Address on which to listen to git requests.")
	authUserFlag      = flag.String("auth.user", "", "Username for basic auth.")
	authPassFlag      = flag.String("auth.pass", "", "Password for basic auth.")
	reposRoot         = flag.String("repos.root", fmt.Sprintf("%s/repos", os.Getenv("HOME")), "The location of the repositories")
	autoInitRepos     = flag.Bool("repos.autoinit", false, "Auto inits repositories on git-push")
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

	config := &gitserver.Config{
		ReposRoot:     *reposRoot,
		AutoInitRepos: *autoInitRepos,
	}

	err := gitserver.Init(config)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on localhost:4000")
	app := gitserver.MiddlewareFunc(handler)
	if hasUserAndPassword() {
		app = authMiddleware(app)
	}
	http.HandleFunc("/", app)
	http.ListenAndServe(*listenAddressFlag, nil)
}
