package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Print network list",
	Long: `Display a network list as table by default 
	as json by using the flag --json or -j`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		networks, err := client.GetNetworkList(ctx)
		if err != nil {
			panic(err)
		}
		var networkinfos [][]string
		if !jsonFlag {
			for _, netw := range networks {
				fill := [][]string{
					{
						netw.Properties.Name,
						netw.Properties.LocationName,
						netw.Properties.CreateTime.String()[:10],
						netw.Properties.Status,
					},
				}
				networkinfos = append(networkinfos, fill...)
			}
			render.Table(out, []string{"network", "location", "createtime", "status"}, networkinfos)
			fmt.Print(out)
		} else {
			fmt.Println(render.AsJSON(networks))
		}
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
