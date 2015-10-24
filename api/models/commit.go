package models

import (
	"bufio"
	"fmt"
	"github.com/dcu/git-http-server/gitserver"
	"strconv"
	"strings"
)

var commitFormat = `--pretty=format:commit:%H%nauthor:%aN%nemail:%aE%ndate:%at%nsubject:%s%nbody:%b`

// Commit stores information about a commit
type Commit struct {
	ID          string
	Author      string
	AuthorEmail string
	Subject     string
	Body        string
	Date        int
}

// ToPublicResponse returns the fields to be send to the client.
func (commit *Commit) ToPublicResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":      commit.ID,
		"author":  commit.Author,
		"email":   commit.AuthorEmail,
		"subject": commit.Subject,
		"body":    commit.Body,
		"date":    commit.Date,
	}
}

// GetCommitsFromBranch gets `count` commits from the given `branch`
func GetCommitsFromBranch(branch string, count int) []*Commit {
	commits := []*Commit{}

	cmd := gitserver.GitCommand{Args: []string{"log", fmt.Sprintf("-%d", count), commitFormat}}
	output, err := cmd.Run(false)
	if err != nil {
		return commits
	}

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
