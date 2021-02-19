package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serverCmdFlags struct {
	forceShutdown    bool
	memory           int
	cores            int
	storageSize      int
	serverName       string
	templateName     string
	hostName         string
	plainPassword    string
	profile          string
	availabilityZone string
	autoRecovery     bool
	includeRelated   bool
	force            bool
}

var (
	serverFlags serverCmdFlags
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Operations on servers",
	Long:  `List, create, or remove servers.`,
}

func serverLsCmdRun(cmd *cobra.Command, args []string) error {
	serverOp := rt.ServerOperator()
	ctx := context.Background()
	out := new(bytes.Buffer)
	servers, err := serverOp.GetServerList(ctx)
	if err != nil {
		return NewError(cmd, "Could not get list of servers", err)
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
		} else {
			render.AsTable(out, heading, rows, renderOpts)
		}
	} else {
		render.AsJSON(out, servers)
	}
	fmt.Print(out)
	return nil
}

var serverLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List servers",
	Long:    `List server objects.`,
	RunE:    serverLsCmdRun,
}

func serverOnCmdRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	serverOp := rt.ServerOperator()
	err := serverOp.StartServer(ctx, args[0])
	if err != nil {
		return NewError(cmd, "Failed starting server", err)
	}
	return nil
}

var serverOnCmd = &cobra.Command{
	Use:   "on ID",
	Short: "Turn server on",
	Args:  cobra.ExactArgs(1),
	RunE:  serverOnCmdRun,
}

func serverOffCmdRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	serverOp := rt.ServerOperator()
	if serverFlags.forceShutdown {
		err := serverOp.StopServer(ctx, args[0])
		if err != nil {
			return NewError(cmd, "Failed stopping server", err)
		}
	} else {
		err := serverOp.ShutdownServer(ctx, args[0])
		if err != nil {
			return NewError(cmd, "Failed shutting down server", err)
		}
	}
	return nil
}

var serverOffCmd = &cobra.Command{
	Use:   "off [flags] ID",
	Short: "Turn server off via ACPI",
	Args:  cobra.ExactArgs(1),
	RunE:  serverOffCmdRun,
}

func serverRmCmdRun(cmd *cobra.Command, args []string) error {
	serverOp := rt.ServerOperator()
	ctx := context.Background()
	id := args[0]
	s, err := serverOp.GetServer(ctx, id)
	if err != nil {
		return NewError(cmd, "Look up server failed", err)
	}
	if serverFlags.force {
		if s.Properties.Power {
			err := serverOp.StopServer(ctx, args[0])
			if err != nil {
				return NewError(cmd, "Failed stopping server", err)
			}
		}
	}

	var storages []gsclient.ServerStorageRelationProperties
	var ipAddrs []gsclient.ServerIPRelationProperties
	if serverFlags.includeRelated {
		out := new(bytes.Buffer)
		storages, err = rt.ServerStorageRelationOperator().GetServerStorageList(ctx, id)
		if err != nil {
			return NewError(cmd, "Could not get related storages", err)
		}
		ipAddrs, err = rt.ServerIPRelationOperator().GetServerIPList(ctx, id)
		if err != nil {
			return NewError(cmd, "Could not get assigned IP addresses", err)
		}

		if !rootFlags.quiet {
			var rows [][]string
			heading := []string{"id", "type", "name"}
			rows = append(rows, []string{
				id,
				"Server",
				s.Properties.Name,
			})
			for _, storage := range storages {
				fill := [][]string{
					{
						storage.ObjectUUID,
						"Storage",
						storage.ObjectName,
					},
				}
				rows = append(rows, fill...)
			}
			for _, addr := range ipAddrs {
				fill := [][]string{
					{
						addr.ObjectUUID,
						fmt.Sprintf("IPv%d address", addr.Family),
						addr.IP,
					},
				}
				rows = append(rows, fill...)
			}

			render.AsTable(out, heading, rows, renderOpts)
			fmt.Print(out)
		}

		if !serverFlags.force {
			msg := "This can destroy your data. "
			if rootFlags.quiet {
				msg += "Re-run with --force to remove"
			} else {
				msg += "Re-run with --force to remove above objects"
			}
			log.Println(msg)
			return nil
		}
	}
	err = serverOp.DeleteServer(ctx, id)
	if err != nil {
		return NewError(cmd, "Deleting server failed", err)
	}
	fmt.Fprintf(os.Stderr, "Removed %s\n", id)

	if serverFlags.includeRelated {
		storageOp := rt.StorageOperator()
		for _, storage := range storages {
			err = storageOp.DeleteStorage(ctx, storage.ObjectUUID)
			if err != nil {
				return NewError(cmd, "Failed removing storage", err)
			}
			fmt.Fprintf(os.Stderr, "Removed %s\n", storage.ObjectUUID)
		}

		ipOp := rt.IPOperator()
		for _, addr := range ipAddrs {
			err = ipOp.DeleteIP(ctx, addr.ObjectUUID)
			if err != nil {
				return NewError(cmd, "Failed removing IP address", err)
			}
			fmt.Fprintf(os.Stderr, "Removed %s\n", addr.ObjectUUID)
		}
	}
	return nil
}

var serverRmCmd = &cobra.Command{
	Use:     "rm [flags] ID",
	Aliases: []string{"remove"},
	Short:   "Remove server",
	Long: `**gscloud server rm** removes an existing server from a project.

With the **--all** option, you can delete all referenced storages and assigned IP addresses, if any. By default, storages and IP addresses are not removed to prevent important data from being deleted.

# EXAMPLES

Remove a server including storages and IP addresses:

	$ gscloud server rm --include-related --force 37d53278-8e5f-47e1-a63f-54513e4b4d53
`,
	Args: cobra.ExactArgs(1),
	RunE: serverRmCmdRun,
}

var serverCreateCmd = &cobra.Command{
	Use:     "create [flags]",
	Example: `gscloud server create --name "My machine" --cores 2 --mem 4 --with-template "My template" --hostname myhost`,
	Short:   "Create server",
	Long: `Create a new server.

# EXAMPLES

Create a server with 25 GB storage from the CentOS 8 template:

	$ gscloud server create \
		--name worker-1 \
		--cores=2 \
		--mem=4 \
		--with-template="CentOS 8 (x86_64)" \
		--storage-size=25 \
		--hostname worker-1

To create a server without any storage just omit --with-template flag:

	$ gscloud server create --name worker-2 --cores=1 --mem=1
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var template gsclient.Template

		serverOp := rt.ServerOperator()
		ctx := context.Background()
		profile, err := toHardwareProfile(serverFlags.profile)
		if err != nil {
			return NewError(cmd, "Cannot create server", err)
		}

		if serverFlags.templateName != "" {
			templateOp := rt.TemplateOperator()
			template, err = templateOp.GetTemplateByName(ctx, serverFlags.templateName)
			if err != nil {
				return NewError(cmd, "Cannot create server", err)
			}
		}

		cleanupServer := false
		server, err := serverOp.CreateServer(ctx, gsclient.ServerCreateRequest{
			Name:            serverFlags.serverName,
			Cores:           serverFlags.cores,
			Memory:          serverFlags.memory,
			HardwareProfile: profile,
			AvailablityZone: serverFlags.availabilityZone,
			AutoRecovery:    &serverFlags.autoRecovery,
		})
		if err != nil {
			return NewError(cmd, "Creating server failed", err)
		}
		cleanupServer = true
		defer func() {
			if cleanupServer {
				err = serverOp.DeleteServer(ctx, server.ObjectUUID)
				if err != nil {
					panic(err)
				}
			}
		}()

		if serverFlags.templateName != "" {
			var password string

			if serverFlags.plainPassword == "" {
				password = generatePassword()
			} else {
				password = serverFlags.plainPassword
			}

			storageOp := rt.StorageOperator()
			storage, err := storageOp.CreateStorage(ctx, gsclient.StorageCreateRequest{
				Name:        string(serverFlags.serverName),
				Capacity:    serverFlags.storageSize,
				StorageType: gsclient.DefaultStorageType,
				Template: &gsclient.StorageTemplate{
					TemplateUUID: template.Properties.ObjectUUID,
					Password:     password,
					PasswordType: gsclient.PlainPasswordType,
					Hostname:     serverFlags.hostName,
				},
			})
			if err != nil {
				return NewError(cmd, "Creating storage failed", err)
			}

			serverStorageOp := rt.ServerStorageRelationOperator()
			err = serverStorageOp.CreateServerStorage(
				ctx,
				server.ObjectUUID,
				gsclient.ServerStorageRelationCreateRequest{
					ObjectUUID: storage.ObjectUUID,
					BootDevice: true,
				})
			if err != nil {
				return NewError(cmd, "Linking storage to server failed", err)
			}
			cleanupServer = false
			fmt.Println("Server created:", server.ObjectUUID)
			fmt.Println("Storage created:", storage.ObjectUUID)
			fmt.Println("Password:", password)
		}

		cleanupServer = false
		return nil
	},
}

var serverSetCmd = &cobra.Command{
	Use:     "set [flags] ID",
	Example: `gscloud server set 37d53278-8e5f-47e1-a63f-54513e4b4d53 --cores 4`,
	Short:   "Update server",
	Long:    `Update properties of an existing server.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return NewError(cmd, "Failed setting property", err)
		}
		return nil
	},
}

var serverAssignCmd = &cobra.Command{
	Use:     "assign ID ADDR",
	Example: `gscloud server assign 37d53278-8e5f-47e1-a63f-54513e4b4d53 2001:db8:0:1::1c8`,
	Short:   "Assign an IP address",
	Long:    `Assign an existing IP address to a server.`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
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
				return err
			}
		} else {
			addrID = args[1]
		}
		err = rt.Client().LinkIP(ctx, serverID, addrID)
		if err != nil {
			return NewError(cmd, "Could not assign IP address", err)
		}
		return nil
	},
}

var serverEventsCmd = &cobra.Command{
	Use:     "events ID",
	Example: `gscloud server events 37d53278-8e5f-47e1-a63f-54513e4b4d53`,
	Short:   "List events",
	Long: `Retrieve event log for given server.
# EXAMPLES

List all events of a server:

	$ gscloud server events 37d53278-8e5f-47e1-a63f-54513e4b4d53

Only list request IDs of a server (in case you need to tell suport what happened):

	$ gscloud server events --quiet 37d53278-8e5f-47e1-a63f-54513e4b4d53

`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverID := args[0]
		ctx := context.Background()
		serverOp := rt.ServerOperator()
		events, err := serverOp.GetServerEventList(ctx, serverID)
		if err != nil {
			return NewError(cmd, "Could not get list of events", err)
		}

		out := new(bytes.Buffer)
		if rootFlags.json {
			render.AsJSON(out, events)

		} else {
			if rootFlags.quiet {
				for _, event := range events {
					fmt.Println(event.Properties.RequestUUID)
				}
			} else {
				var rows [][]string
				heading := []string{
					"time", "request id", "request type", "details", "initiator",
				}
				for _, event := range events {
					fill := [][]string{
						{
							event.Properties.Timestamp.Local().Format(time.RFC3339),
							event.Properties.RequestUUID,
							event.Properties.RequestType,
							event.Properties.Change,
							event.Properties.Initiator,
						},
					}
					rows = append(rows, fill...)
				}
				render.AsTable(out, heading, rows, renderOpts)
			}
		}
		fmt.Print(out)

		return nil
	},
}

func init() {
	serverOffCmd.Flags().BoolVarP(&serverFlags.forceShutdown, "force", "f", false, "Force shutdown (no ACPI)")

	serverCreateCmd.Flags().IntVar(&serverFlags.memory, "mem", 1, "Memory (GB)")
	serverCreateCmd.Flags().IntVar(&serverFlags.cores, "cores", 1, "No. of cores")
	serverCreateCmd.Flags().IntVar(&serverFlags.storageSize, "storage-size", 10, "Storage capacity (GB)")
	serverCreateCmd.Flags().StringVar(&serverFlags.serverName, "name", "", "Name of the server")
	serverCreateCmd.Flags().StringVar(&serverFlags.templateName, "with-template", "", "Name of template to use")
	serverCreateCmd.Flags().StringVar(&serverFlags.hostName, "hostname", "", "Hostname")
	serverCreateCmd.Flags().StringVar(&serverFlags.plainPassword, "password", "", "Plain-text password")
	serverCreateCmd.Flags().MarkDeprecated("password", "a password will be created automatically")
	serverCreateCmd.Flags().Lookup("password").Hidden = true
	serverCreateCmd.Flags().StringVar(&serverFlags.profile, "profile", "q35", "Hardware profile")
	serverCreateCmd.Flags().StringVar(&serverFlags.availabilityZone, "availability-zone", "", "Availability zone. One of \"a\", \"b\", \"c\". Default \"\"")
	serverCreateCmd.Flags().BoolVar(&serverFlags.autoRecovery, "auto-recovery", true, "Whether to restart in case of errors")

	serverSetCmd.Flags().IntVar(&serverFlags.memory, "mem", 0, "Memory (GB)")
	serverSetCmd.Flags().IntVar(&serverFlags.cores, "cores", 0, "No. of cores")
	serverSetCmd.Flags().StringVar(&serverFlags.serverName, "name", "", "Name of the server")

	serverRmCmd.Flags().BoolVarP(&serverFlags.includeRelated, "include-related", "i", false, "Remove all objects currently related to this server, not just the server")
	serverRmCmd.Flags().BoolVarP(&serverFlags.force, "force", "f", false, "Force a destructive operation")

	serverCmd.AddCommand(serverLsCmd, serverOnCmd, serverOffCmd, serverRmCmd, serverCreateCmd, serverSetCmd, serverAssignCmd, serverEventsCmd)
	rootCmd.AddCommand(serverCmd)
}

func generatePassword() string {
	res, err := password.Generate(12, 6, 2, false, false)
	if err != nil {
		panic(err)
	}
	return res
}

func toHardwareProfile(val string) (gsclient.ServerHardwareProfile, error) {
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
		return "", fmt.Errorf("Not a valid profile: %s", val)
	}
	return prof, nil
}
