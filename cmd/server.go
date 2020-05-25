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
	Short: "Print server list",
	Long: `Display server list as table by default 
	as json by using the flag --json or -j`,
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
						strconv.FormatInt(int64(server.Properties.Cores), 10),
						strconv.FormatInt(int64(server.Properties.Memory), 10),
						server.Properties.Status,
					},
				}
				serverinfos = append(serverinfos, fill...)

			}
			render.Table(out, []string{"server-name", "cores", "memory", "status"}, serverinfos)
			fmt.Print(out)
		} else {
			fmt.Println(render.AsJSON(servers))

		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
