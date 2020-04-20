package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GitCommit gets its value through Makefile
var GitCommit string

// Version gets its value through Makefile
var Version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns last Git Commit and Version",
	Long: `gscloud version displays latest Git commit SHA and Version number.
	
For example:
./gscloud version
Version:        0.2.0-beta
Git commit:     66a9631ed5c1516d34ca305d4432149b67675cd0`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:\t%s\nGit commit:\t%s\n", Version, GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
