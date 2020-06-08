package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Print network list",
	Long:  `Print all networks information`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		networks, err := client.GetNetworkList(ctx)
		if err != nil {
			panic(err)
		}
		var networkinfos [][]string
		if !jsonFlag {
			heading := []string{"name", "location", "createtime", "status", "uuid"}
			for _, netw := range networks {
				fill := [][]string{
					{
						netw.Properties.Name,
						netw.Properties.LocationName,
						netw.Properties.CreateTime.String()[:10],
						netw.Properties.Status,
						netw.Properties.ObjectUUID,
					},
				}
				networkinfos = append(networkinfos, fill...)

			}
			if idFlag {
				rowsToDisplay = len(heading)
			}
			render.Table(out, heading[:rowsToDisplay], networkinfos)

		} else {
			render.AsJSON(out, networks)
		}
		fmt.Print(out)
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
