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

// serverOperator is used for testing purpose,
// we can mock data return from the gsclient via interface.
type serverOperator interface {
	GetServerList(ctx context.Context) ([]gsclient.Server, error)
	StartServer(ctx context.Context, id string) error
	StopServer(ctx context.Context, id string) error
	ShutdownServer(ctx context.Context, id string) error
	DeleteServer(ctx context.Context, id string) error
	//CreateServer(ctx context.Context, request gsclient.ServerCreateRequest) error
}

// Server action enums
const (
	serverListAction = iota
	serverStartAction
	serverStopAction
	serverShutdownAction
	serverDeleteAction
	serverCreateAction
)

var forceFlag bool

// produceServerCmdRunFunc takes an instance of a struct that implements `serverOperator`
// returns a `cmdRunFunc`
func produceServerCmdRunFunc(o serverOperator, action int) cmdRunFunc {
	switch action {
	case serverListAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			out := new(bytes.Buffer)
			servers, err := o.GetServerList(ctx)
			if err != nil {
				log.Fatalf("Couldn't get server list: %s", err)
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
			err := o.StartServer(ctx, args[0])
			if err != nil {
				log.Fatalf("Failed starting server: %s", err)
			}
		}

	case serverStopAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if forceFlag {
				err := o.StopServer(ctx, args[0])
				if err != nil {
					log.Fatalf("Failed stopping server: %s", err)
				}
			} else {
				err := o.ShutdownServer(ctx, args[0])
				if err != nil {
					log.Fatalf("Failed shutting down server: %s", err)
				}
			}
		}
	case serverDeleteAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := o.DeleteServer(ctx, args[0])
			if err != nil {
				log.Fatalf("Removing Server failed: %s", err)
			}
		}
	case serverCreateAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			c, m := args[1], args[2]
			cor, mem := 0, 0
			if n, err := strconv.Atoi(c); err == nil {
				cor = n
			}
			if n, err := strconv.Atoi(m); err == nil {
				mem = n
			}
			cServer, err := client.CreateServer(ctx, gsclient.ServerCreateRequest{
				Name:   args[0],
				Cores:  cor,
				Memory: mem,
			})
			if err != nil {
				log.Fatal("Create server has failed with error", err)
			}
			log.WithFields(log.Fields{
				"server_uuid": cServer.ObjectUUID,
			}).Info("Server successfully created")
		}
	default:
	}
	return nil
}

func initServerCmd() {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Operations on servers",
		Long:  `List, create, or remove servers.`,
	}

	var serverLsCmd = &cobra.Command{
		Use:     "ls [flags]",
		Aliases: []string{"list"},
		Short:   "List servers",
		Long:    `List server objects.`,
		Run:     produceServerCmdRunFunc(client, serverListAction),
	}

	var serverOnCmd = &cobra.Command{
		Use:   "on [ID]",
		Short: "Turn server on",
		Args:  cobra.ExactArgs(1),
		Run:   produceServerCmdRunFunc(client, serverStartAction),
	}

	var serverOffCmd = &cobra.Command{
		Use:   "off [flags] [ID]",
		Short: "Turn server off via ACPI",
		Args:  cobra.ExactArgs(1),
		Run:   produceServerCmdRunFunc(client, serverStopAction),
	}
	var removeCmd = &cobra.Command{
		Use:     "rm [flags] [ID]",
		Aliases: []string{"remove"},
		Short:   "Remove server",
		Long:    `Remove an existing server.`,
		Args:    cobra.ExactArgs(1),
		Run:     produceServerCmdRunFunc(client, serverDeleteAction),
	}
	var createCmd = &cobra.Command{
		Use:     "create [Name] [Cores] [Memory]",
		Aliases: []string{"create"},
		Short:   "Create server",
		Long:    `Create a new server.`,
		Args:    cobra.ExactArgs(3),
		Run:     produceServerCmdRunFunc(client, serverCreateAction),
	}
	serverOffCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "Force shutdown (no ACPI)")
	serverCmd.AddCommand(serverLsCmd, serverOnCmd, serverOffCmd, removeCmd, createCmd)
	rootCmd.AddCommand(serverCmd)
}
