package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Print server list",
	Long:  `Print all server information`,
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
				status := "off"
				if server.Properties.Power {
					status = "on"
				}
				fill := [][]string{
					{
						server.Properties.Name,
						strconv.FormatInt(int64(server.Properties.Cores), 10),
						strconv.FormatInt(int64(server.Properties.Memory), 10),
						status,
					},
				}
				serverinfos = append(serverinfos, fill...)

			}
			render.Table(out, []string{"server", "core", "mem", "power"}, serverinfos)
			fmt.Print(out)
		} else {
			render.AsJSON(servers)

		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
