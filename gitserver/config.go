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
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Path      string `yaml:"path"`
	DisableUI bool   `yaml:"disable"`
}

// AuthyConfig is the config for the Authy 2FA service
type AuthyConfig struct {
	APIKey string `yaml:"api_key"`
	UserID string `yaml:"user_id"`
}

// Config stores the config of the git server
type Config struct {
	Host       string       `yaml:"host"`
	EnableCORS bool         `yaml:"cors"`
	Repos      *ReposConfig `yaml:"repos"`
	UI         *UIConfig    `yaml:"ui"`
	Authy      *AuthyConfig `yaml:"authy"`
}

// HasAuth returns whether the auth is configured or not.
func (config *Config) HasAuth() bool {
	return config.UI.Username != "" && config.UI.Password != ""
}

// HasAuthy returns true if Authy is configured.
func (config *Config) HasAuthy() bool {
	return config.Authy != nil && config.Authy.APIKey != "" && config.Authy.UserID != ""
}

// WriteToPath writes the config to the given filePath
func (config *Config) WriteToPath(filePath string) {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		panic(err)
	}
}

// WriteSampleConfig writes a default config to the given path
func WriteSampleConfig(path string) {
	config := Config{
		Host:       ":4000",
		EnableCORS: true,
		Repos: &ReposConfig{
			AutoInit: true,
			Path:     "/tmp/repos",
		},
		UI: &UIConfig{
			Username:  "admin",
			Password:  "admin",
			Path:      "/repos",
			DisableUI: false,
		},
		Authy: &AuthyConfig{
			APIKey: "",
			UserID: "",
		},
	}
	config.WriteToPath(path)
}

// LoadConfig loads the config from the given path
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
