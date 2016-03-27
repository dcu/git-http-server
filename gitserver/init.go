package gitserver

import (
	"os"
)

var gServerConfig *Config

// Init initializes the git server.
func Init(config *Config) error {
	gServerConfig = config

	err := os.MkdirAll(gServerConfig.Repos.Path, 0700)

	return err
}

// ReposRoot returns the directory where repositories are stored.
func ReposRoot() string {
	return gServerConfig.Repos.Path
}
