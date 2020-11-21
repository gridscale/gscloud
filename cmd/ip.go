package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ipCmdFlags struct {
	v4         bool
	v6         bool
	failover   bool
	name       string
	reverseDNS string
}

var (
	ipFlags ipCmdFlags
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Operations on IP addresses",
	Long:  `List, add, or remove IP address objects.`,
}

var ipLsCmd = &cobra.Command{
	Use:     "ls [-4|-6]",
	Aliases: []string{"list"},
	Short:   "List IP addresses",
	Long:    `List IP address objects.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if ipFlags.v4 && ipFlags.v6 {
			log.Fatal("No family selected")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ipOp := rt.IPOperator()
		ctx := context.Background()
		ipAddresses, err := ipOp.GetIPList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get list of IPs: %s", err)
		}
		var rows [][]string
		out := new(bytes.Buffer)
		if !rootFlags.json {
			heading := []string{"IP", "assigned", "failover", "family", "reverse DNS", "ID"}
			for _, addr := range ipAddresses {
				if ipFlags.v4 && addr.Properties.Family == 6 {
					continue
				}

				if ipFlags.v6 && addr.Properties.Family == 4 {
					continue
				}

				isFailover := "no"
				if addr.Properties.Failover {
					isFailover = "yes"
				}
				isAssigned := "free"
				relations := addr.Properties.Relations
				if len(relations.Servers) > 0 || len(relations.Loadbalancers) > 0 {
					isAssigned = "assigned"
				}
				properties := [][]string{
					{
						addr.Properties.IP,
						isAssigned,
						isFailover,
						fmt.Sprintf("v%d", addr.Properties.Family),
						addr.Properties.ReverseDNS,
						addr.Properties.ObjectUUID,
					},
				}
				rows = append(rows, properties...)
			}
			if rootFlags.quiet {
				for _, row := range rows {
					fmt.Println(row[5])
				}
				return
			}
			render.AsTable(out, heading, rows, renderOpts)
		} else {
			render.AsJSON(out, ipAddresses)
		}
		fmt.Print(out)
	},
}

var ipRmCmd = &cobra.Command{
	Use:     "rm [flags] ID|ADDRESS",
	Aliases: []string{"remove"},
	Short:   "Delete an IP address",
	Long: `Remove an existing IP address object by ID or address.

# EXAMPLES

Delete by ID:

    $ gscloud ip rm 71d85c9d-6fdd-404b-a821-1d94c2050a6e

Delete by address:

    $ gscloud ip rm 2a06:2380:2:1::24

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id string
		var err error
		ctx := context.Background()
		ipOp := rt.IPOperator()
		address := net.ParseIP(args[0])
		if address != nil {
			id, err = idForAddress(ctx, address, ipOp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			id = args[0]
		}
		err = ipOp.DeleteIP(ctx, id)
		if err != nil {
			log.Fatalf("Removing IP failed: %s", err)
		}
	},
}

var ipSetCmd = &cobra.Command{
	Use:     "set [flags] ID|ADDRESS",
	Example: `gscloud ip set ID|ADDRESS --reverse-dns example.com`,
	Short:   "Update IP address properties",
	Long: `Update properties of an existing IP address.

# EXAMPLES

Set PTR entry and name on an existing IP:

    $ gscloud ip set 2a06:2380:2:1::85 --name test --reverse-dns example.com
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id string
		var err error
		address := net.ParseIP(args[0])
		ctx := context.Background()
		ipOp := rt.IPOperator()
		if address != nil {
			id, err = idForAddress(ctx, address, ipOp)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			id = args[0]
		}
		updateReq := gsclient.IPUpdateRequest{}
		if ipFlags.failover {
			updateReq.Failover = true
		}
		if len(ipFlags.name) > 0 {
			updateReq.Name = ipFlags.name
		}
		if len(ipFlags.reverseDNS) > 0 {
			updateReq.ReverseDNS = ipFlags.reverseDNS
		}
		err = ipOp.UpdateIP(
			ctx,
			id,
			updateReq)
		if err != nil {
			log.Fatalf("Failed: %s", err)
		}
	},
}

var ipAddCmd = &cobra.Command{
	Use:     "add -4|-6 [flags]",
	Aliases: []string{"create"},
	Example: `gscloud ip add -6`,
	Short:   "Create a new IP address",
	Long: `Create a new IP address object.

# EXAMPLES

Create a new IPv6 address with a PTR entry and a name:

    $ gscloud ip add -6 --name test --reverse-dns=example.com

Create a new IPv4 address:

    $ gscloud ip add -4

`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if ipFlags.v4 && ipFlags.v6 {
			log.Fatal("No family given. Use either -4 or -6")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		family := gsclient.IPv4Type
		if ipFlags.v6 {
			family = gsclient.IPv6Type
		}
		ipOp := rt.IPOperator()
		ctx := context.Background()
		ipAddress, err := ipOp.CreateIP(ctx, gsclient.IPCreateRequest{
			Family:     family,
			Failover:   ipFlags.failover,
			Name:       ipFlags.name,
			ReverseDNS: ipFlags.reverseDNS,
		})
		if err != nil {
			log.Fatalf("Adding IPv%d address failed: %s", family, err)
		}
		fmt.Println("IP added:", ipAddress.IP)
	},
}

func init() {
	ipLsCmd.PersistentFlags().BoolVarP(&ipFlags.v4, "v4", "4", false, "IPv4 only")
	ipLsCmd.PersistentFlags().BoolVarP(&ipFlags.v6, "v6", "6", false, "IPv6 only")

	ipAddCmd.PersistentFlags().BoolVarP(&ipFlags.v4, "v4", "4", false, "Add a new IPv4 address")
	ipAddCmd.PersistentFlags().BoolVarP(&ipFlags.v6, "v6", "6", false, "Add a new IPv6 address")
	ipAddCmd.PersistentFlags().StringVarP(&ipFlags.name, "name", "n", "", "Optional name of the IP address being created. Can be omitted")
	ipAddCmd.PersistentFlags().BoolVarP(&ipFlags.failover, "failover", "", false, "Enable failover. If given, IP is no longer available for DHCP and cannot be assigned")
	ipAddCmd.PersistentFlags().StringVarP(&ipFlags.reverseDNS, "reverse-dns", "", "", "Optional reverse DNS entry for the IP address")

	ipSetCmd.PersistentFlags().StringVarP(&ipFlags.name, "name", "n", "", "Change name of the IP address")
	ipSetCmd.PersistentFlags().BoolVarP(&ipFlags.failover, "failover", "", false, "Enable failover")
	ipSetCmd.PersistentFlags().StringVarP(&ipFlags.reverseDNS, "reverse-dns", "", "", "Set reverse DNS entry")

	ipCmd.AddCommand(ipLsCmd, ipRmCmd, ipSetCmd, ipAddCmd)
	rootCmd.AddCommand(ipCmd)
}

func idForAddress(ctx context.Context, addr net.IP, op gsclient.IPOperator) (string, error) {
	ipAddresses, err := op.GetIPList(ctx)
	if err != nil {
		return "", err
	}
	for _, obj := range ipAddresses {
		ip := net.ParseIP(obj.Properties.IP)
		if ip != nil && ip.Equal(addr) {
			return obj.Properties.ObjectUUID, nil
		}
	}
	return "", fmt.Errorf("No such IP %s", addr)
}
