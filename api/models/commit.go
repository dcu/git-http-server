package models

import ()

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
