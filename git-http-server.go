package main

import (
	"flag"
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"github.com/dcu/http-einhorn"
	"github.com/gin-gonic/gin"
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
			w.Header().Set("WWW-Authenticate", "Basic realm=\"git-server\"")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			next(w, r)
		}
	}
}

func hasUserAndPassword() bool {
	return *authUserFlag != "" && *authPassFlag != ""
}

func handler(c *gin.Context) {
	c.String(200, "nothing to see here\n")
}

func gitserverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedRoute := gitserver.MatchRoute(c.Request)
		if parsedRoute != nil {
			parsedRoute.Dispatch(c.Writer, c.Request)
		} else {
			c.Next()
		}
	}
}

func startHTTP() {
	log.Printf("Starting server on %s", *listenAddressFlag)

	router := gin.Default()
	if hasUserAndPassword() {
		router.Use(gin.BasicAuth(gin.Accounts{*authUserFlag: *authPassFlag}))
	}
	router.Use(gitserverHandler())
	router.GET("/", handler)

	if einhorn.IsRunning() {
		einhorn.Start(router, 0)
	} else {
		router.Run(*listenAddressFlag)
	}
}

func parseOptsAndBuildConfig() *gitserver.Config {
	flag.Parse()

	config := &gitserver.Config{
		ReposRoot:     *reposRoot,
		AutoInitRepos: *autoInitRepos,
	}

	return config
}

func main() {
	config := parseOptsAndBuildConfig()

	err := gitserver.Init(config)
	if err != nil {
		panic(err)
	}

	startHTTP()
}
