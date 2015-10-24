package api

import (
	"github.com/dcu/git-http-server/api/models"
	"github.com/dcu/git-http-server/gitserver"
	"github.com/gin-gonic/gin"
)

func listRepos(c *gin.Context) {
	repos := models.FindAllRepositories(gitserver.ReposRoot())

	response := gin.H{}

	repositoriesInfo := []gin.H{}
	for _, repo := range repos {
		repositoriesInfo = append(repositoriesInfo, repo.ToPublicResponse())
	}
	response["items"] = repositoriesInfo
	response["items_count"] = len(repositoriesInfo)

	c.JSON(200, response)
}
