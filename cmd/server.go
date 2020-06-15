package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

var forceFlag bool

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
			heading := []string{"name", "core", "mem", "power", "id"}
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
						server.Properties.ObjectUUID,
					},
				}
				serverinfos = append(serverinfos, fill...)
			}
			if idFlag {
				rowsToDisplay = len(heading)
			}
			render.Table(out, heading[:rowsToDisplay], serverinfos)
		} else {
			render.AsJSON(out, servers)
		}
		fmt.Print(out)
	},
}

var onCmd, offCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn server on",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client.StartServer(ctx, args[0])
	},
}, &cobra.Command{
	Use:   "off",
	Short: "Turn server off",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if forceFlag {
			client.StopServer(ctx, args[0])
		}
		client.ShutdownServer(ctx, args[0])
	},
}

func init() {
	serverCmd.AddCommand(onCmd, offCmd)
	serverCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "Force shutdown")
	rootCmd.AddCommand(serverCmd)
}
