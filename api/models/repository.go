package models

import (
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	testGitFiles = []string{"HEAD", "info", "objects", "refs"}
)

// Repository has the information about a repository.
type Repository struct {
	Path string
}

// NewRepository creates an instance of a Repository.
func NewRepository(path string) *Repository {
	return &Repository{
		Path: path,
	}
}

// ToPublicResponse returns the fields to expose in the public response
func (repo *Repository) ToPublicResponse() map[string]interface{} {
	return map[string]interface{}{
		"name":        repo.Name(),
		"last_commit": repo.LastCommit().ToPublicResponse(),
	}
}

// LastCommit returns the last commit on this repository
func (repo *Repository) LastCommit() *Commit {
	return GetCommitsFromBranch("HEAD", 1)[0]
}

// Name returns the name of the repository
func (repo *Repository) Name() string {
	name := strings.TrimPrefix(repo.Path, gitserver.ReposRoot())
	name = strings.TrimSuffix(name, ".git")
	return name
}

// ReadmeFile returns the contents of the readme file.
func (repo *Repository) ReadmeFile(branch string) string {
	filePath := fmt.Sprintf("%s:%s", resolveBranch(branch), "README.md")

	cmd := gitserver.GitCommand{Args: []string{"show", filePath}}
	output := cmd.RunAndGetOutput()

	return string(output)
}

// FindAllRepositories find all repositories
func FindAllRepositories(dir string) []*Repository {
	repos := []*Repository{}

	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		absPath := filepath.Join(dir, file.Name())
		if isGitDir(absPath, file) {
			repos = append(repos, NewRepository(absPath))
		} else if file.IsDir() {
			childrenRepos := FindAllRepositories(absPath)
			repos = append(repos, childrenRepos...)
		}
	}

	return repos
}

func isGitDir(absPath string, fileInfo os.FileInfo) bool {
	if !fileInfo.IsDir() {
		return false
	}

	for _, gitFile := range testGitFiles {
		gitDirPath := filepath.Join(absPath, gitFile)
		_, err := os.Stat(gitDirPath)

		if err != nil {
			return false
		}
	}

	return true
}

func resolveBranch(branch string) string {
	if branch == "" {
		return "HEAD"
	}

	return branch
}
