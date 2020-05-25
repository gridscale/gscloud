package cmd

import (
	"fmt"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/kardianos/osext"
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

type accountEntry struct {
	Name   string `yaml:"name"`
	UserID string `yaml:"userId"`
	Token  string `yaml:"token"`
}
type cliConfig struct {
	Accounts []accountEntry `yaml:"accounts"`
}

func cliPath() string {
	filePath, err := osext.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return filePath
}

// configPath construct platform specific path to the configuration file.
// - on Linux: $XDG_CONFIG_HOME or $HOME/.config
// - on macOS: $HOME/Library/Application Support
// - on Windows: %APPDATA% or "C:\\Users\\%USER%\\AppData\\Roaming"
func cliConfigPath() string {
	path := viper.ConfigFileUsed()
	if path == "" {
		path = configdir.LocalConfig("gscloud") + "/config.yaml"
		viper.SetConfigFile(path)
	}
	return path
}

func cliCachePath() string {
	return configdir.LocalCache("gscloud")
}

func newCliClient(account string) *gsclient.Client {
	var ac accountEntry

	cliConf := &cliConfig{}
	err := viper.Unmarshal(cliConf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, a := range cliConf.Accounts {
		if account == a.Name {
			ac = a
			break
		}
	}

	// clientConf := &clientConfig{
	// apiURL:     defaultAPIURL,
	// userUUID:   ac.UserID,
	// userToken:  ac.Token,
	// userAgent:  "gscloud",
	// httpClient: http.DefaultClient,
	// }
	// return newClient(clientConf)
	config := gsclient.DefaultConfiguration(
		ac.UserID,
		ac.Token,
	)
	return gsclient.NewClient(config)

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
