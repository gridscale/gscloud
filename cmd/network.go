package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// networkOperator is used for testing purpose,
// we can mock data return from the gsclient via interface.
type networkOperator interface {
	GetNetworkList(ctx context.Context) ([]gsclient.Network, error)
}

// produceNetworkCmdRunFunc takes an instance of a struct that implements `networkOperator`
// returns a `cmdRunFunc`
func produceNetworkCmdRunFunc(o networkOperator) cmdRunFunc {
	return func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		networks, err := o.GetNetworkList(ctx)
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
	}
}

func initNetworkCmd() {
	var networkCmd, lsCmd = &cobra.Command{
		Use:   "network",
		Short: "Print network list",
		Long:  `Print all networks information`,
		Run:   produceNetworkCmdRunFunc(client),
	}, &cobra.Command{
		Use:   "ls",
		Short: "List networks",
		Args:  cobra.MaximumNArgs(1),
		Run:   produceNetworkCmdRunFunc(client),
	}
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(lsCmd)
}
