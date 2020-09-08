package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gridscale/gscloud/runtime"
	"github.com/gridscale/gscloud/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var makeConfigCmd = &cobra.Command{
	Use:   "make-config",
	Short: "Create a new configuration file",
	Long: fmt.Sprintf(`Create a new and possibly almost empty configuration file overwriting an
existing one if it exists. Prints the path to the newly created file to
stdout.

Examples:

Create a new configuration file at the default configuration path
(%s):

gscloud make-config

Create a new configuration file at a specified path:

gscloud --config ~/myconfig.yaml make-config

`, runtime.ConfigPath()),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := runtime.ConfigPath()

		if !utils.FileExists(filePath) {
			err := os.MkdirAll(filepath.Dir(filePath), os.FileMode(0700))
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(filePath, emptyConfig(), 0644)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Written: %s\n", filePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(makeConfigCmd)
}

// emptyConfig creates a new config YAML with a 'default' account.
func emptyConfig() []byte {
	defaultAccount := runtime.AccountEntry{
		Name:   "default",
		UserID: "",
		Token:  "",
		URL:    defaultAPIURL,
	}
	c := runtime.Config{
		Accounts: []runtime.AccountEntry{defaultAccount},
	}
	out, _ := yaml.Marshal(&c)
	return out
}
