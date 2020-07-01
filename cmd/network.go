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
	Short: "Print network list",
	Long:  `Print all networks information`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		networks, err := client.GetNetworkList(ctx)
		if err != nil {
			log.Error("Couldn't get Networkinfo", err)
			return
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

func init() {
	rootCmd.AddCommand(networkCmd)
}
