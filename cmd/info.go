package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the info",
	Long:  `Print information belongs to gscloud accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := runtime.ParseConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		for _, account := range conf.Accounts {
			accountName := rt.Account()
			// get info of the current account
			if account.Name == accountName {
				// Get info about primitive resources
				client := rt.Client()
				servers, err := client.GetServerList(context.Background())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(3)
				}
				storages, err := client.GetStorageList(context.Background())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(3)
				}
				ipAddrs, err := client.GetIPList(context.Background())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(3)
				}
				paasServices, err := client.GetPaaSServiceList(context.Background())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(3)
				}
				fmt.Printf(
					"Account: %s\nUserID: %s\nToken: %s\nURL: %s\n",
					account.Name, account.UserID, account.Token, account.URL)
				fmt.Printf(
					"No. of servers: %d\nNo. of storages: %d\nNo. of ip addresses: %d\nNo. of platform services: %d\n",
					len(servers), len(storages), len(ipAddrs), len(paasServices))
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Hide some global persistent flags here that don't make sense on 'version'
	origHelpFunc := versionCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "info" || (cmd.HasParent() && cmd.Parent().Name() == "info") {
			cmd.Flags().MarkHidden("account")
			cmd.Flags().MarkHidden("config")
			cmd.Flags().MarkHidden("json")
			cmd.Flags().MarkHidden("quiet")
		}
		origHelpFunc(cmd, args)
	})
}
