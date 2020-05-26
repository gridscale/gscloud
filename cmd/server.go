package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Print server list.",
	Long:  `Print all server information.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		out := new(bytes.Buffer)
		servers, err := client.GetServerList(ctx)
		if err != nil {
			panic(err)
		}
		var serverinfos [][]string
		if !jsonFlag {
			for _, server := range servers {

				fill := [][]string{
					{
						server.Properties.Name,
						string("Core/s " + strconv.FormatInt(int64(server.Properties.Cores), 10) + " RAM " + strconv.FormatInt(int64(server.Properties.Memory), 10)),
						strconv.FormatBool(server.Properties.Power),
						strconv.FormatInt(int64(server.Properties.CurrentPrice), 10),
					},
				}
				serverinfos = append(serverinfos, fill...)

			}
			render.Table(out, []string{"server", "specifications", "power", "currentprice"}, serverinfos)
			fmt.Print(out)
		} else {
			fmt.Println(render.AsJSON(servers))

		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
