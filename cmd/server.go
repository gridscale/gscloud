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

var (
	forceShutdown bool
	memory        int
	cores         int
	storage       int
	serverName    string
	template      string
	hostName      string
	plainPassword string
)

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
	Run: func(cmd *cobra.Command, args []string) {
		serverOp := rt.ServerOperator()
		ctx := context.Background()
		out := new(bytes.Buffer)
		servers, err := serverOp.GetServerList(ctx)
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
	},
}

var serverOnCmd = &cobra.Command{
	Use:   "on [ID]",
	Short: "Turn server on",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverOp := rt.ServerOperator()
		err := serverOp.StartServer(ctx, args[0])
		if err != nil {
			log.Fatalf("Failed starting server: %s", err)
		}
	},
}

var serverOffCmd = &cobra.Command{
	Use:   "off [flags] [ID]",
	Short: "Turn server off via ACPI",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverOp := rt.ServerOperator()
		if forceShutdown {
			err := serverOp.StopServer(ctx, args[0])
			if err != nil {
				log.Fatalf("Failed stopping server: %s", err)
			}
		} else {
			err := serverOp.ShutdownServer(ctx, args[0])
			if err != nil {
				log.Fatalf("Failed shutting down server: %s", err)
			}
		}
	},
}

var serverRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove server",
	Long:    `Remove an existing server.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverOp := rt.ServerOperator()
		ctx := context.Background()
		err := serverOp.DeleteServer(ctx, args[0])
		if err != nil {
			log.Fatalf("Removing server failed: %s", err)
		}
	},
}

var serverCreateCmd = &cobra.Command{
	Use:     "create [flags]",
	Example: `./gscloud server create --name "My machine" --cores 2 --mem 4 --with-template "My template" --password mysecret --hostname myhost`,
	Short:   "Create server",
	Long:    `Create a new server.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverOp := rt.ServerOperator()
		ctx := context.Background()
		cServer, err := serverOp.CreateServer(ctx, gsclient.ServerCreateRequest{
			Name:   serverName,
			Cores:  cores,
			Memory: memory,
		})
		if err != nil {
			log.Fatalf("Creating server failed: %s", err)
		}
		fmt.Println("Server created:", cServer.ObjectUUID)

		if template != "" {
			template, _ := serverOp.GetTemplateByName(ctx, template)
			cStorage, err := serverOp.CreateStorage(ctx, gsclient.StorageCreateRequest{
				Name:        string(serverName),
				Capacity:    storage,
				StorageType: gsclient.DefaultStorageType,
				Template: &gsclient.StorageTemplate{
					TemplateUUID: template.Properties.ObjectUUID,
					Password:     plainPassword,
					PasswordType: gsclient.PlainPasswordType,
					Hostname:     hostName,
				},
			})
			serverOp.CreateServerStorage(
				ctx,
				cServer.ObjectUUID,
				gsclient.ServerStorageRelationCreateRequest{
					ObjectUUID: cStorage.ObjectUUID,
					BootDevice: true,
				})
			if err != nil {
				log.Fatalf("Create storage failed: %s", err)
			}
			fmt.Println("Storage created:", cStorage.ObjectUUID)
		}
	},
}

func init() {
	serverOffCmd.PersistentFlags().BoolVarP(&forceShutdown, "force", "f", false, "Force shutdown (no ACPI)")

	serverCreateCmd.PersistentFlags().IntVar(&memory, "mem", 1, "Memory (GB)")
	serverCreateCmd.PersistentFlags().IntVar(&cores, "cores", 1, "No. of cores")
	serverCreateCmd.PersistentFlags().IntVar(&storage, "storage-size", 10, "Storage capacity (GB)")
	serverCreateCmd.PersistentFlags().StringVar(&serverName, "name", "", "Name of the server")
	serverCreateCmd.PersistentFlags().StringVar(&template, "with-template", "", "Name of template to use")
	serverCreateCmd.PersistentFlags().StringVar(&hostName, "hostname", "", "Hostname")
	serverCreateCmd.PersistentFlags().StringVar(&plainPassword, "password", "", "Plain-text password")

	serverCmd.AddCommand(serverLsCmd, serverOnCmd, serverOffCmd, serverRmCmd, serverCreateCmd)
	rootCmd.AddCommand(serverCmd)
}
