package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
)

type infoJSONOutput struct {
	runtime.AccountEntry
	ServerCount  int `json:"server_count"`
	StorageCount int `json:"storage_count"`
	IPAddrCount  int `json:"ip_address_count"`
	PaaSCount    int `json:"paas_service_count"`
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the info",
	Long:  `Print information belongs to gscloud accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := runtime.ParseConfig()
		if err != nil {
			return NewError(cmd, "Could parse configuration", err)
		}
		for _, account := range conf.Accounts {
			accountName := rt.Account()
			// get info of the current account
			if account.Name == accountName {
				// Print auth info
				if !rootFlags.json {
					out := new(bytes.Buffer)
					heading := []string{"Account", "UserID", "Token", "URL"}
					fill := [][]string{
						{
							account.Name,
							account.UserID,
							account.Token,
							account.URL,
						},
					}
					var rows [][]string
					rows = append(rows, fill...)
					render.AsTable(out, heading, rows, renderOpts)
					fmt.Print(out)
					fmt.Printf("\nGetting infomation about available resources...\n\n\n")
				}

				// Get info about primitive resources
				client := rt.Client()
				servers, err := client.GetServerList(context.Background())
				if err != nil {
					return NewError(cmd, "Could not get servers' information", err)
				}
				storages, err := client.GetStorageList(context.Background())
				if err != nil {
					return NewError(cmd, "Could not get storages' information", err)
				}
				ipAddrs, err := client.GetIPList(context.Background())
				if err != nil {
					return NewError(cmd, "Could not get ip addresses' information", err)
				}
				paasServices, err := client.GetPaaSServiceList(context.Background())
				if err != nil {
					return NewError(cmd, "Could not get PaaS services' information", err)
				}
				out := new(bytes.Buffer)
				if !rootFlags.json {
					heading := []string{"No. of servers", "No. of storages", "No. of ip addresses", "No. of platform services"}
					fill := [][]string{
						{
							strconv.Itoa(len(servers)),
							strconv.Itoa(len(storages)),
							strconv.Itoa(len(ipAddrs)),
							strconv.Itoa(len(paasServices)),
						},
					}
					var rows [][]string
					rows = append(rows, fill...)
					render.AsTable(out, heading, rows, renderOpts)
				} else {
					jsonOutput := infoJSONOutput{
						AccountEntry: account,
						ServerCount:  len(servers),
						StorageCount: len(storages),
						IPAddrCount:  len(ipAddrs),
						PaaSCount:    len(paasServices),
					}
					render.AsJSON(out, jsonOutput)
				}
				fmt.Print(out)
			}
		}
		return nil
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
			// cmd.Flags().MarkHidden("json")
			cmd.Flags().MarkHidden("quiet")
		}
		origHelpFunc(cmd, args)
	})
}
