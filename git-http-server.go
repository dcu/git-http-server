package main

import (
	"flag"
	"fmt"
	"github.com/dcu/git-http-server/api"
	"github.com/dcu/git-http-server/gitserver"
	"github.com/dcu/http-einhorn"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

var (
	listenAddressFlag = flag.String("web.listen-address", ":4000", "Address on which to listen to git requests.")
	authUserFlag      = flag.String("auth.user", "", "Username for basic auth.")
	authPassFlag      = flag.String("auth.pass", "", "Password for basic auth.")
	reposRoot         = flag.String("repos.root", fmt.Sprintf("%s/repos", os.Getenv("HOME")), "The location of the repositories.")
	autoInitRepos     = flag.Bool("repos.autoinit", false, "Auto inits repositories on git-push.")
	reposPath         = flag.String("web.ui_path", "/repos", "HTTP path where repos UI can be found.")
	disableUI         = flag.Bool("web.disable_ui", false, "Disables web UI")
	disableCors       = flag.Bool("web.disable_cors", false, "Disables Cross-Origin Resource Sharing")
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
	if !*disableUI {
		c.Redirect(http.StatusTemporaryRedirect, *reposPath)
	} else {

		c.String(200, "nothing to see here\n")
	}
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

func corsHandler() gin.HandlerFunc {
	handler := cors.Default()
	return func(c *gin.Context) {
		handler.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}

func startHTTP() {
	log.Printf("Starting server on %s", *listenAddressFlag)

	router := gin.Default()
	if hasUserAndPassword() {
		router.Use(gin.BasicAuth(gin.Accounts{*authUserFlag: *authPassFlag}))
	}
	router.Use(gitserverHandler())

	if !*disableCors {
		router.Use(corsHandler())
	}

	if !*disableUI {
		router.StaticFS(*reposPath, http.Dir("./public/"))
	}
	router.GET("/", handler)

	api.SetupRouter(router)

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
