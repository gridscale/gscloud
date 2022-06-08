package runtime

import (
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

// AccountEntry represents a single account in the config file.
type AccountEntry struct {
	Name   string `yaml:"name" json:"name"`
	UserID string `yaml:"userId" json:"userId"`
	Token  string `yaml:"token" json:"token"`
	URL    string `yaml:"url" json:"url"`
}

// Config are all configuration settings parsed from a configuration file.
type Config struct {
	Accounts []AccountEntry `yaml:"accounts"`
}

// ConfigPath construct platform specific path to the configuration file.
// - on Linux: $XDG_CONFIG_HOME or $HOME/.config
// - on macOS: $HOME/Library/Application Support
// - on Windows: %APPDATA% or "C:\\Users\\%USER%\\AppData\\Roaming"
func ConfigPath() string {
	path := viper.ConfigFileUsed()
	if path == "" {
		path = configdir.LocalConfig("gscloud") + "/config.yaml"
		viper.SetConfigFile(path)
	}
	return path
}

// ConfigPathWithoutUser same as ConfigPath but with environment variables not expanded.
func ConfigPathWithoutUser() string {
	return localConfig + "/gscloud/config.yaml"
}

// ParseConfig parses viper config file.
func ParseConfig() (*Config, error) {
	conf := Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
