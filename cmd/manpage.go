package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manpageCmd = &cobra.Command{
	Use:   "manpage [PATH]",
	Short: "Create man-pages for gscloud",
	Long: `Build and write man-pages to given path.
Example:

Create a new set of section 1 man-pages in /usr/local:

gscloud manpage /usr/local/share/man/man1

This will overwrite any existing man-page created previously.
`,
	Example: `gscloud manpage /path/to/man-pages`,
	Run: func(cmd *cobra.Command, args []string) {
		header := &doc.GenManHeader{
			Title:   "GSCLOUD",
			Section: "1",
			Source:  " ",
		}
		err := doc.GenManTree(rootCmd, header, args[0])
		if err != nil {
			log.Fatalf("Couldn't create man-pages: %s", err)
		}
		fmt.Println("Written to:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(manpageCmd)
}
