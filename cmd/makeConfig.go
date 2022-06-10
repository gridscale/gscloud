package cmd

import (
	"fmt"

	"github.com/gridscale/gscloud/runtime"
	"github.com/gridscale/gscloud/utils"
	"github.com/spf13/cobra"
)

var makeConfigCmd = &cobra.Command{
	Use:   "make-config",
	Short: "Create a new configuration file",
	Long: fmt.Sprintf(`Create a new and possibly almost empty configuration file not overwriting an existing one if it exists. Prints the path to the newly created file to stdout.

# EXAMPLES

Create a new configuration file at the default configuration path (%s):

    $ gscloud make-config

Create a new configuration file at a specified path:

    $ gscloud --config /tmp/myconfig.yaml make-config

`, runtime.ConfigPathWithoutUser()),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := runtime.ConfigPath()

		if rootFlags.configFile != "" {
			filePath = rootFlags.configFile
		}

		if !utils.FileExists(filePath) {
			runtime.WriteConfig(&runtime.Config{Projects: []runtime.ProjectEntry{{URL: defaultAPIURL}}}, filePath)

			fmt.Printf("Written: %s\n", filePath)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(makeConfigCmd)
}
