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
	DeleteNetwork(ctx context.Context, id string) error
}

// Network action enums
const (
	networkListAction = iota
	networkDeleteAction
)

// produceNetworkCmdRunFunc takes an instance of a struct that implements `networkOperator`
// returns a `cmdRunFunc`
func produceNetworkCmdRunFunc(o networkOperator, action int) cmdRunFunc {
	switch action {
	case networkListAction:
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
	case networkDeleteAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := o.DeleteNetwork(ctx, args[0])
			if err != nil {
				log.Fatalf("Removing Network failed: %s", err)
			}
		}

	default:
	}
	return nil
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
		Run:     produceNetworkCmdRunFunc(client, networkListAction),
	}
	var removeCmd = &cobra.Command{
		Use:     "rm [flags] [ID]",
		Aliases: []string{"remove"},
		Short:   "Remove Network",
		Long:    `Remove an existing Network.`,
		Args:    cobra.ExactArgs(1),
		Run:     produceNetworkCmdRunFunc(client, networkDeleteAction),
	}

	networkCmd.AddCommand(networkLsCmd, removeCmd)
	rootCmd.AddCommand(networkCmd)
}
