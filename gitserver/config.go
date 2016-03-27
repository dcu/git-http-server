package gitserver

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ReposConfig is the configuration of the repositories
type ReposConfig struct {
	Path     string `yaml:"path"`
	AutoInit bool   `yaml:"autoinit"`
}

// UIConfig is the config of the UI
type UIConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Path      string `yaml:"path"`
	DisableUI bool   `yaml:"disable"`
}

// Config stores the config of the git server
type Config struct {
	Host       string       `yaml:"host"`
	EnableCORS bool         `yaml:"cors"`
	Repos      *ReposConfig `yaml:"repos"`
	UI         *UIConfig    `yaml:"ui"`
}

// HasAuth returns whether the auth is configured or not.
func (config *Config) HasAuth() bool {
	return config.UI.Username != "" && config.UI.Password != ""
}

func LoadConfig(path string) *Config {
	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(data, config)

	return config
}
