package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	Long:    `List network objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		networkOps := rt.NetworkOperator()
		networks, err := networkOps.GetNetworkList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get network list: %s", err)
		}
		var networkinfos [][]string
		if !jsonFlag {
			heading := []string{"id", "name", "location", "createtime", "status"}
			for _, netw := range networks {
				fill := [][]string{
					{
						netw.Properties.ObjectUUID,
						netw.Properties.Name,
						netw.Properties.LocationName,
						netw.Properties.CreateTime.String()[:10],
						netw.Properties.Status,
					},
				}
				networkinfos = append(networkinfos, fill...)

			}
			render.Table(out, heading[:], networkinfos)
			if quietFlag {
				for _, info := range networkinfos {
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

var networkRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove Network",
	Long:    `Remove an existing Network.`,
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
	networkCmd.AddCommand(networkLsCmd, networkRmCmd)
	rootCmd.AddCommand(networkCmd)
}
