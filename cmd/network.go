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
	}
}

func initNetworkCmd() {
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
		Run:     produceNetworkCmdRunFunc(client),
	}

	networkCmd.AddCommand(networkLsCmd)
	rootCmd.AddCommand(networkCmd)
}
