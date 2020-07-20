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
var servernameFlag, templateFlag, passwordFlag, hostnameFlag string
var cpuFlag, memFlag, storageFlag int

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
			cServer, err := client.CreateServer(ctx, gsclient.ServerCreateRequest{
				Name:   servernameFlag,
				Cores:  cpuFlag,
				Memory: memFlag,
			})
			if err != nil {
				log.Fatal("Create server has failed with error", err)
			}
			log.WithFields(log.Fields{
				"server_uuid": cServer.ObjectUUID,
			}).Info("Server successfully created")

			if templateFlag != "" {
				template, _ := client.GetTemplateByName(ctx, templateFlag)
				cStorage, err := client.CreateStorage(ctx, gsclient.StorageCreateRequest{
					Name:        string(servernameFlag),
					Capacity:    storageFlag,
					StorageType: gsclient.DefaultStorageType,
					Template: &gsclient.StorageTemplate{
						TemplateUUID: template.Properties.ObjectUUID,
						Password:     passwordFlag,
						PasswordType: gsclient.PlainPasswordType,
						Hostname:     hostnameFlag,
					},
					Labels: []string{"label"},
				})

				client.CreateServerStorage(
					ctx,
					cServer.ObjectUUID,
					gsclient.ServerStorageRelationCreateRequest{
						ObjectUUID: cStorage.ObjectUUID,
						BootDevice: true,
					})

				if err != nil {
					log.Error("Create storage has failed with error", err)
					return
				}
				log.WithFields(log.Fields{
					"storage_uuid": cStorage.ObjectUUID,
				}).Info("Storage successfully created")
				produceStorageCmdRunFunc(client, storageCreateAction)

			}

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
		Use:     "create [flags]",
		Example: `./gscloud server create --name "Server Name" --cpu 6 --mem 4 --with-template "Debian 10" --password "p4ssw0rd" --hostname "gridscale"`,
		Aliases: []string{"create"},
		Short:   "Create server",
		Long:    `Create a new server.`,
		Run:     produceServerCmdRunFunc(client, serverCreateAction),
	}

	// Flags
	serverOffCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "Force shutdown (no ACPI)")
	createCmd.PersistentFlags().StringVarP(&servernameFlag, "name", "n", "default_server", "Define servername")
	createCmd.PersistentFlags().IntVarP(&cpuFlag, "cpu", "c", 1, "Define cpu cores")
	createCmd.PersistentFlags().IntVarP(&memFlag, "mem", "m", 1, "Define Memory")
	createCmd.PersistentFlags().StringVarP(&passwordFlag, "password", "p", "pass", "Define password")
	createCmd.PersistentFlags().StringVarP(&hostnameFlag, "hostname", "H", "no_hostname", "Define hostname")
	createCmd.PersistentFlags().StringVarP(&templateFlag, "with-template", "t", "Debian 10", "Define Operating System")
	createCmd.PersistentFlags().IntVarP(&storageFlag, "storage-size", "s", 10, "Storage size")
	//
	serverCmd.AddCommand(serverLsCmd, serverOnCmd, serverOffCmd, removeCmd, createCmd)
	rootCmd.AddCommand(serverCmd)

}
