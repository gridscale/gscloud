package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log"
)

var manpagePath string

// Manpage action enum
const (
	manpageCreateAction = iota
)

// produceManpageCmdRunFunc takes an instance of a struct that implements `manpageOperator`
// returns a `cmdRunFunc`
func produceManpageCmdRunFunc(action int) cmdRunFunc {
	switch action {
	case manpageCreateAction:
		return func(cmd *cobra.Command, args []string) {

			header := &doc.GenManHeader{
				Title:   "GSCLOUD",
				Section: "1",
				Source:  "Copyright (c) 2020 gridscale GmbH",
			}
			fmt.Print(manpagePath)
			err := doc.GenManTree(rootCmd, header, manpagePath)
			if err != nil {
				log.Fatalf("Couldn't create Manpage: %s", err)
			}
		}
	default:
	}
	return nil
}

func initManpageCmd() {
	var manpageCreateCmd = &cobra.Command{
		Use:     "manpage [flags]",
		Short:   "Add Manpage",
		Long:    `Create manpage in given path.`,
		Example: `./gscloud manpage --path "/path/to/your/manpages`,
		Run:     produceManpageCmdRunFunc(manpageCreateAction),
	}
	manpageCreateCmd.PersistentFlags().StringVarP(&manpagePath, "path", "p", "", "Export Manpage to given path")

	rootCmd.AddCommand(manpageCreateCmd)
}
