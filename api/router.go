package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter setup api routes
func SetupRouter(router *gin.Engine) {
	router.GET("/repositories", listRepos)
	router.GET("/repositories/*path", showRepo)
}
