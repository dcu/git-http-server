package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dcu/git-http-server/api"
	"github.com/dcu/git-http-server/gitserver"
	"github.com/dcu/http-einhorn"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

var (
	gConfig           *gitserver.Config
	defaultConfigPath = "config.yml"
)

func handler(c *gin.Context) {
	if !gConfig.UI.DisableUI {
		c.Redirect(http.StatusTemporaryRedirect, gConfig.Repos.Path)
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
	log.Printf("Starting server on %s", gConfig.Host)

	router := gin.Default()
	if gConfig.HasAuth() {
		router.Use(gin.BasicAuth(gin.Accounts{gConfig.UI.Username: gConfig.UI.Password}))
	}
	router.Use(gitserverHandler())

	if gConfig.EnableCORS {
		router.Use(corsHandler())
	}

	if !gConfig.UI.DisableUI {
		router.StaticFS(gConfig.Repos.Path, http.Dir("./public/"))
	}
	router.GET("/", handler)

	api.SetupRouter(router)

	if einhorn.IsRunning() {
		einhorn.Start(router, 0)
	} else {
		router.Run(gConfig.Host)
	}
}

func parseOptsAndBuildConfig() (*gitserver.Config, error) {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
			fmt.Printf("Writing sample config to %s\n", defaultConfigPath)
			gitserver.WriteSampleConfig(defaultConfigPath)
		}

		flag.Usage()
	}

	return gitserver.LoadConfig(args[0])
}

func main() {
	config, err := parseOptsAndBuildConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %s\n", err.Error())
		return
	}

	err = gitserver.Init(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config server: %s\n", err.Error())
		return
	}

	gConfig = config
	startHTTP()
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: git-http-server <config.yml>\n")
		flag.PrintDefaults()

		os.Exit(1)
	}
}
