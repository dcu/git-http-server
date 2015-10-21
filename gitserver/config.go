package gitserver

// Config stores the config of the git server
type Config struct {
	ReposRoot     string "repos_root"
	AutoInitRepos bool   "repos_autoinit"
}
