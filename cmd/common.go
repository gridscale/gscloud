package cmd

import "github.com/spf13/cobra"

type cmdRunFunc = func(cmd *cobra.Command, args []string)
