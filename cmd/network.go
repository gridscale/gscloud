package cmd

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type networkCmdFlags struct {
	networkName string
}

var (
	networkFlags networkCmdFlags
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Operations on networks",
	Long:  `List, create, or remove networks.`,
}

var networkLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List networks",
	Long:    `List networks.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		networkOps := rt.NetworkOperator()
		networks, err := networkOps.GetNetworkList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get network list: %s", err)
		}
		var rows [][]string
		if !rootFlags.json {
			heading := []string{"id", "name", "location", "changed", "status"}
			for _, network := range networks {
				fill := [][]string{
					{
						network.Properties.ObjectUUID,
						network.Properties.Name,
						network.Properties.LocationName,
						network.Properties.ChangeTime.Local().Format(time.RFC3339),
						network.Properties.Status,
					},
				}
				rows = append(rows, fill...)

			}
			render.AsTable(out, heading, rows, renderOpts)
			if rootFlags.quiet {
				for _, info := range rows {
					fmt.Println(info[0])
				}
				return
			}

		} else {
			render.AsJSON(out, networks)
		}
		fmt.Print(out)
	},
}

var networkCreateCmd = &cobra.Command{
	Use:     "create [flags]",
	Example: `gscloud network create --name myNetwork`,
	Short:   "Create network",
	Long: `Create a new network.

# EXAMPLES

Create a network:

	$ gscloud network create 

`,
	Run: func(cmd *cobra.Command, args []string) {
		networkOp := rt.NetworkOperator()
		ctx := context.Background()
		network, err := networkOp.CreateNetwork(ctx, gsclient.NetworkCreateRequest{
			Name: networkFlags.networkName,
		})

		if err != nil {
			log.Fatalf("Creating network failed: %s", err)
		}
		fmt.Println("Network created:", network.ObjectUUID)
	},
}

var networkRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove network",
	Long:    `Remove an existing network.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		networkOps := rt.NetworkOperator()
		err := networkOps.DeleteNetwork(ctx, args[0])
		if err != nil {
			log.Fatalf("Removing network failed: %s", err)
		}
	},
}

func init() {
	networkCreateCmd.Flags().StringVar(&networkFlags.networkName, "name", "", "Name of the network")

	networkCmd.AddCommand(networkLsCmd, networkRmCmd, networkCreateCmd)
	rootCmd.AddCommand(networkCmd)
}
