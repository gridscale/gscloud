package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serverCmdFlags struct {
	forceShutdown bool
	memory        int
	cores         int
	storage       int
	serverName    string
	template      string
	hostName      string
	plainPassword string
	profile       string
}

var (
	serverFlags serverCmdFlags
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Operations on servers",
	Long:  `List, create, or remove servers.`,
}

func serverLsCmdRun(cmd *cobra.Command, args []string) {
	serverOp := rt.ServerOperator()
	ctx := context.Background()
	out := new(bytes.Buffer)
	servers, err := serverOp.GetServerList(ctx)
	if err != nil {
		log.Fatalf("Couldn't get server list: %s", err)
	}
	var rows [][]string
	if !rootFlags.json {
		heading := []string{"id", "name", "core", "mem", "changed", "power"}
		for _, server := range servers {
			power := "off"
			if server.Properties.Power {
				power = "on"
			}
			fill := [][]string{
				{
					server.Properties.ObjectUUID,
					server.Properties.Name,
					strconv.FormatInt(int64(server.Properties.Cores), 10),
					strconv.FormatInt(int64(server.Properties.Memory), 10),
					server.Properties.ChangeTime.Local().Format(time.RFC3339),
					power,
				},
			}
			rows = append(rows, fill...)
		}
		if rootFlags.quiet {
			for _, info := range rows {
				fmt.Println(info[0])
			}
			return
		}
		render.AsTable(out, heading, rows, renderOpts)
	} else {
		render.AsJSON(out, servers)
	}
	fmt.Print(out)
}

var serverLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List servers",
	Long:    `List server objects.`,
	Run:     serverLsCmdRun,
}

func serverOnCmdRun(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	serverOp := rt.ServerOperator()
	err := serverOp.StartServer(ctx, args[0])
	if err != nil {
		log.Fatalf("Failed starting server: %s", err)
	}
}

var serverOnCmd = &cobra.Command{
	Use:   "on ID",
	Short: "Turn server on",
	Args:  cobra.ExactArgs(1),
	Run:   serverOnCmdRun,
}

func serverOffCmdRun(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	serverOp := rt.ServerOperator()
	if serverFlags.forceShutdown {
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
}

var serverOffCmd = &cobra.Command{
	Use:   "off [flags] ID",
	Short: "Turn server off via ACPI",
	Args:  cobra.ExactArgs(1),
	Run:   serverOffCmdRun,
}

func serverRmCmdRun(cmd *cobra.Command, args []string) {
	serverOp := rt.ServerOperator()
	ctx := context.Background()
	err := serverOp.DeleteServer(ctx, args[0])
	if err != nil {
		log.Fatalf("Removing server failed: %s", err)
	}
}

var serverRmCmd = &cobra.Command{
	Use:     "rm [flags] ID",
	Aliases: []string{"remove"},
	Short:   "Remove server",
	Long:    `Remove an existing server.`,
	Args:    cobra.ExactArgs(1),
	Run:     serverRmCmdRun,
}

var serverCreateCmd = &cobra.Command{
	Use:     "create [flags]",
	Example: `gscloud server create --name "My machine" --cores 2 --mem 4 --with-template "My template" --password mysecret --hostname myhost`,
	Short:   "Create server",
	Long:    `Create a new server.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverOp := rt.ServerOperator()
		ctx := context.Background()
		profile := toHardwareProfile(serverFlags.profile)
		cServer, err := serverOp.CreateServer(ctx, gsclient.ServerCreateRequest{
			Name:            serverFlags.serverName,
			Cores:           serverFlags.cores,
			Memory:          serverFlags.memory,
			HardwareProfile: profile,
		})
		if err != nil {
			log.Fatalf("Creating server failed: %s", err)
		}
		fmt.Println("Server created:", cServer.ObjectUUID)

		if serverFlags.template != "" {
			templateOp := rt.TemplateOperator()
			template, _ := templateOp.GetTemplateByName(ctx, serverFlags.template)

			storageOp := rt.StorageOperator()
			cStorage, err := storageOp.CreateStorage(ctx, gsclient.StorageCreateRequest{
				Name:        string(serverFlags.serverName),
				Capacity:    serverFlags.storage,
				StorageType: gsclient.DefaultStorageType,
				Template: &gsclient.StorageTemplate{
					TemplateUUID: template.Properties.ObjectUUID,
					Password:     serverFlags.plainPassword,
					PasswordType: gsclient.PlainPasswordType,
					Hostname:     serverFlags.hostName,
				},
			})

			serverStorageOp := rt.ServerStorageRelationOperator()
			serverStorageOp.CreateServerStorage(
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

var serverSetCmd = &cobra.Command{
	Use:     "set [flags] ID",
	Example: `gscloud server set 37d53278-8e5f-47e1-a63f-54513e4b4d53 --cores 4`,
	Short:   "Update server",
	Long:    `Update properties of an existing server.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverOp := rt.ServerOperator()
		ctx := context.Background()
		err := serverOp.UpdateServer(
			ctx,
			args[0],
			gsclient.ServerUpdateRequest{
				Cores:  serverFlags.cores,
				Memory: serverFlags.memory,
				Name:   serverFlags.serverName,
			})
		if err != nil {
			log.Fatalf("Failed: %s", err)
		}
	},
}

var serverAssignCmd = &cobra.Command{
	Use:     "assign ID ADDR",
	Example: `gscloud server assign 37d53278-8e5f-47e1-a63f-54513e4b4d53 2001:db8:0:1::1c8`,
	Short:   "Assign an IP address",
	Long:    `Assign an existing IP address to a server.`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var serverID string
		var addrID string
		var err error

		serverID = args[0]
		ctx := context.Background()
		ipOp := rt.IPOperator()
		addr := net.ParseIP(args[1])
		if addr != nil {
			addrID, err = idForAddress(ctx, addr, ipOp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			addrID = args[1]
		}
		err = rt.Client().LinkIP(ctx, serverID, addrID)
		if err != nil {
			log.Fatalf("Failed: %s", err)
		}
	},
}

func init() {
	serverOffCmd.PersistentFlags().BoolVarP(&serverFlags.forceShutdown, "force", "f", false, "Force shutdown (no ACPI)")

	serverCreateCmd.PersistentFlags().IntVar(&serverFlags.memory, "mem", 1, "Memory (GB)")
	serverCreateCmd.PersistentFlags().IntVar(&serverFlags.cores, "cores", 1, "No. of cores")
	serverCreateCmd.PersistentFlags().IntVar(&serverFlags.storage, "storage-size", 10, "Storage capacity (GB)")
	serverCreateCmd.PersistentFlags().StringVar(&serverFlags.serverName, "name", "", "Name of the server")
	serverCreateCmd.PersistentFlags().StringVar(&serverFlags.template, "with-template", "", "Name of template to use")
	serverCreateCmd.PersistentFlags().StringVar(&serverFlags.hostName, "hostname", "", "Hostname")
	serverCreateCmd.PersistentFlags().StringVar(&serverFlags.plainPassword, "password", "", "Plain-text password")
	serverCreateCmd.PersistentFlags().StringVar(&serverFlags.profile, "profile", "q35", "Hardware profile")

	serverSetCmd.PersistentFlags().IntVar(&serverFlags.memory, "mem", 0, "Memory (GB)")
	serverSetCmd.PersistentFlags().IntVar(&serverFlags.cores, "cores", 0, "No. of cores")
	serverSetCmd.PersistentFlags().StringVar(&serverFlags.serverName, "name", "", "Name of the server")

	serverCmd.AddCommand(serverLsCmd, serverOnCmd, serverOffCmd, serverRmCmd, serverCreateCmd, serverSetCmd, serverAssignCmd)
	rootCmd.AddCommand(serverCmd)
}

func toHardwareProfile(val string) gsclient.ServerHardwareProfile {
	var prof gsclient.ServerHardwareProfile
	switch val {
	case "default":
		prof = gsclient.DefaultServerHardware

	case "nested":
		prof = gsclient.NestedServerHardware

	case "legacy":
		prof = gsclient.LegacyServerHardware

	case "cisco_csr":
		prof = gsclient.CiscoCSRServerHardware

	case "sophos_utm":
		prof = gsclient.SophosUTMServerHardware

	case "f5_bigip":
		prof = gsclient.F5BigipServerHardware

	case "q35":
		prof = gsclient.Q35ServerHardware

	case "q35_nested":
		prof = gsclient.Q35NestedServerHardware

	default:
		log.Fatal("Not a valid profile")
	}
	return prof
}
