package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sshKeysOperator is used for testing purpose,
// we can mock data return from the gsclient via interface.
type serverOperator interface {
	GetServerList(ctx context.Context) ([]gsclient.Server, error)
	StartServer(ctx context.Context, id string) error
	StopServer(ctx context.Context, id string) error
	ShutdownServer(ctx context.Context, id string) error
}

// Server action enums
const (
	serverMainAction = iota
	serverStartAction
	serverStopAction
	serverShutdownAction
)

var forceFlag bool

// produceServerCmdRunFunc takes an instance of a struct that implements `serverOperator`
// returns a `cmdRunFunc`
func produceServerCmdRunFunc(o serverOperator, action int) cmdRunFunc {
	switch action {
	case serverMainAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			out := new(bytes.Buffer)
			servers, err := o.GetServerList(ctx)
			if err != nil {
				log.Error("Couldn't get Serverinfo", err)
				return
			}
			var serverinfos [][]string
			if !jsonFlag {
				heading := []string{"id", "name", "core", "mem", "power"}
				for _, server := range servers {
					status := "off"
					if server.Properties.Power {
						status = "on"
					}
					fill := [][]string{
						{
							server.Properties.ObjectUUID,
							server.Properties.Name,
							strconv.FormatInt(int64(server.Properties.Cores), 10),
							strconv.FormatInt(int64(server.Properties.Memory), 10),
							status,
						},
					}
					serverinfos = append(serverinfos, fill...)
				}
				if quietFlag {
					for _, info := range serverinfos {
						fmt.Println(info[0])
					}
					return
				}
				render.Table(out, heading[:], serverinfos)
			} else {
				render.AsJSON(out, servers)
			}
			fmt.Print(out)
		}

	case serverStartAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			o.StartServer(ctx, args[0])
		}

	case serverStopAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if forceFlag {
				o.StopServer(ctx, args[0])
			}
			o.ShutdownServer(ctx, args[0])
		}

	default:
	}
	return nil
}

func initServerCmd() {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Print server list",
		Long:  `Print all server information`,
		Run:   produceServerCmdRunFunc(client, serverMainAction),
	}

	var onCmd, offCmd = &cobra.Command{
		Use:   "on",
		Short: "Turn server on",
		Args:  cobra.MinimumNArgs(1),
		Run:   produceServerCmdRunFunc(client, serverStartAction),
	}, &cobra.Command{
		Use:   "off",
		Short: "Turn server off",
		Args:  cobra.MinimumNArgs(1),
		Run:   produceServerCmdRunFunc(client, serverStopAction),
	}
	serverCmd.AddCommand(onCmd, offCmd)
	serverCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "Force shutdown")
	rootCmd.AddCommand(serverCmd)
}
