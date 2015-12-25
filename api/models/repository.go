package models

import (
	"bufio"
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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
	commits := repo.GetCommitsFromBranch("HEAD", 1)
	if len(commits) > 0 {
		return commits[0]
	}

	return nil
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

	cmd := repo.GitCommand([]string{"show", filePath})
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

// GetCommitsFromBranch gets `count` commits from the given `branch`
func (repo *Repository) GetCommitsFromBranch(branch string, count int) []*Commit {
	commits := []*Commit{}

	output := repo.runGitLogCommand(branch, count)
	scanner := bufio.NewScanner(output)
	var lastCommit *Commit

	for scanner.Scan() {
		line := scanner.Text()
		keyAndValue := strings.SplitN(line, ":", 2)

		switch keyAndValue[0] {
		case "commit":
			{
				lastCommit = &Commit{ID: keyAndValue[1]}
				commits = append(commits, lastCommit)
			}
		case "date":
			{
				date, err := strconv.Atoi(keyAndValue[1])
				if err == nil {
					lastCommit.Date = date
				}
			}
		case "author":
			{
				lastCommit.Author = keyAndValue[1]
			}
		case "email":
			{
				lastCommit.AuthorEmail = keyAndValue[1]
			}
		case "subject":
			{
				lastCommit.Subject = keyAndValue[1]
			}
		case "body":
			{
				lastCommit.Body = keyAndValue[1]
			}
		default:
			{
				lastCommit.Body += "\n" + keyAndValue[0]
			}
		}
	}

	return commits
}

// GitCommand returns a git command based on the repo settings.
func (repo *Repository) GitCommand(args []string) *gitserver.GitCommand {
	argsForRepo := []string{"--git-dir", repo.Path}
	argsForRepo = append(argsForRepo, args...)

	return &gitserver.GitCommand{Args: argsForRepo}
}

func (repo *Repository) runGitLogCommand(branch string, count int) io.ReadCloser {
	cmd := repo.GitCommand([]string{"log", resolveBranch(branch), fmt.Sprintf("-%d", count), commitFormat})
	output, err := cmd.Run(false)
	if err != nil {
		return nil
	}

	return output
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
