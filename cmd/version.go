package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GitCommit value set by a linker flag
var GitCommit string

// Version value set by a linker flag
var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `Print version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:\t%s\nGit commit:\t%s\n", Version, GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Hide some global persistent flags here that don't make sense on 'version'
	origHelpFunc := versionCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "version" || (cmd.HasParent() && cmd.Parent().Name() == "version") {
			cmd.Flags().MarkHidden("project")
			cmd.Flags().MarkHidden("config")
			cmd.Flags().MarkHidden("json")
			cmd.Flags().MarkHidden("quiet")
		}
		origHelpFunc(cmd, args)
	})
}
