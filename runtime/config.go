package runtime

import (
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

// AccountEntry represents a single account in the config file.
type AccountEntry struct {
	Name   string `yaml:"name"`
	UserID string `yaml:"userId"`
	Token  string `yaml:"token"`
	URL    string `yaml:"url"`
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

// CommandWithoutConfig return true if current command does not need a config file.
// Called from within a cobra initializer function. Unfortunately there is no
// way of getting the current command from an cobra initializer so we scan the
// command line again.
func CommandWithoutConfig(cmdLine []string) bool {
	var noConfigNeeded = []string{
		"make-config", "version", "manpage", "completion",
	}
	for _, cmd := range noConfigNeeded {
		if contains(cmdLine, cmd) {
			return true
		}
	}
	return false
}

// contains tests whether string e is in slice s.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
