package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type accountEntry struct {
	Name   string `yaml:"name"`
	UserID string `yaml:"user_id"`
	Token  string `yaml:"token"`
	URL    string `yaml:"url"`
}
type cliConfig struct {
	Accounts []accountEntry `yaml:"accounts"`
}

func cliPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return filepath.Join(dir, os.Args[0])
}

func (c *cliConfig) fetchCliConfig() *cliConfig {
	yamlFile, _ := ioutil.ReadFile(cfgFile)
	err := yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return c
}

func newCliClient(account string) *gsclient {
	var cc cliConfig
	var ac accountEntry
	cliConfig := cc.fetchCliConfig()
	for _, a := range cliConfig.Accounts {
		if account == a.Name {
			ac = a
			break
		}
	}

	config := &config{
		apiURL:     defaultAPIURL,
		userUUID:   ac.UserID,
		userToken:  ac.Token,
		userAgent:  "gscloud",
		httpClient: http.DefaultClient,
	}
	return newClient(config)
}
