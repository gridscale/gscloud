package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh]",
	Short:                 "Generate completion script",
	Long:                  `Generate Tab-Completion script for bash or zsh shell`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
			break
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
			break
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
