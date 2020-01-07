package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type accountEntry struct {
	Name   string `yaml:"name"`
	UserID string `yaml:"userId"`
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

func newCliClient(account string) *gsclient {
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

	clientConf := &clientConfig{
		apiURL:     defaultAPIURL,
		userUUID:   ac.UserID,
		userToken:  ac.Token,
		userAgent:  "gscloud",
		httpClient: http.DefaultClient,
	}
	return newClient(clientConf)
}
