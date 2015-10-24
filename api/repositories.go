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

func showRepo(c *gin.Context) {
	path, _ := c.Params.Get("path")
	absPath, err := gitserver.AbsoluteRepoPath(path)
	if err != nil {
		c.JSON(400, gin.H{"message": "Bad request"})
		return
	}
	repository := models.NewRepository(absPath)

	response := gin.H{
		"repository": repository.ToPublicResponse(),
		"readme":     repository.ReadmeFile("HEAD"),
	}

	c.JSON(200, response)
}
