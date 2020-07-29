package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log"
)

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
			err := doc.GenManTree(rootCmd, header, args[0])
			if err != nil {
				log.Fatalf("Couldn't create Manpage: %s", err)
			}
			fmt.Println("Manpages created:", args[0])
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
		Example: `./gscloud manpage /path/to/manpages`,
		Run:     produceManpageCmdRunFunc(manpageCreateAction),
	}
	rootCmd.AddCommand(manpageCreateCmd)
}
