package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	oldConfigPath = "gscloud"
	configPath    = "gridscale"
)

// AccountEntry represents a single account in the config file.
type ProjectEntry struct {
	Name   string `yaml:"name" json:"name"`
	UserID string `yaml:"userId" json:"userId"`
	Token  string `yaml:"token" json:"token"`
	URL    string `yaml:"url" json:"url"`
}

// Config are all configuration settings parsed from a configuration file.
type Config struct {
	Projects []ProjectEntry `yaml:"projects"`
}

// OldConfig are all configuration settings parsed from an old configuration file
type OldConfig struct {
	Accounts []ProjectEntry `yaml:"accounts"`
}

// ConfigPath constructs the platform specific path to the configuration file.
// - on Linux: $XDG_CONFIG_HOME or $HOME/.config
// - on macOS: $HOME/Library/Application Support
// - on Windows: %APPDATA% or "C:\\Users\\%USER%\\AppData\\Roaming"
func ConfigPath() string {
	p := filepath.Join(configdir.LocalConfig(), configPath)
	return p
}

func OldConfigPath() string {
	p := filepath.Join(configdir.LocalConfig(), oldConfigPath)
	return p
}

// ConfigPathWithoutUser is the same as ConfigPath but with environment variables not expanded.
func ConfigPathWithoutUser() string {
	return localConfig + configPath
}

func OldConfigPathWithoutUser() string {
	return localConfig + oldConfigPath
}

// ParseConfig parses viper config file.
func ParseConfig() (*Config, error) {
	conf := Config{}
	err := viper.Unmarshal(&conf)

	if err != nil {
		return nil, err
	}

	if conf.Projects == nil {
		oldConf := OldConfig{}
		viper.Unmarshal(&oldConf)

		conf.Projects = oldConf.Accounts
	}

	return &conf, nil
}

func WriteConfig(conf *Config, filePath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), os.FileMode(0700))
	if err != nil {
		return err
	}

	c, _ := yaml.Marshal(conf)

	err = ioutil.WriteFile(filePath, c, 0644)
	if err != nil {
		return err
	}

	return nil
}
