package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GitCommit value set by a linker flag
var GitCommit string

// Version value set by a linker flag
var Version string

func versionCmdRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Version:\t%s\nGit commit:\t%s\n", Version, GitCommit)
}

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  `Display gscloud version information.`,
		Run:   versionCmdRun,
	}
	rootCmd.AddCommand(versionCmd)

	// Hide some global persistent flags here that don't make sense on 'version'
	origHelpFunc := versionCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "version" || (cmd.HasParent() && cmd.Parent().Name() == "version") {
			cmd.Flags().MarkHidden("account")
			cmd.Flags().MarkHidden("config")
		}
		origHelpFunc(cmd, args)
	})
}
