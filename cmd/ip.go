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

var (
	v4         bool
	v6         bool
	failover   bool
	name       string
	reverseDNS string
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
		if v4 && v6 {
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
		if !jsonFlag {
			heading := []string{"IP", "assigned", "failover", "family", "reverse DNS", "ID"}
			for _, addr := range ipAddresses {
				if v4 && addr.Properties.Family == 6 {
					continue
				}

				if v6 && addr.Properties.Family == 4 {
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
			if quietFlag {
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

Examples:

Delete by ID:

gscloud ip rm 71d85c9d-6fdd-404b-a821-1d94c2050a6e

Delete by address:

gscloud ip rm 2a06:2380:2:1::24

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

var ipAddCmd = &cobra.Command{
	Use:     "add -4|-6 [flags]",
	Aliases: []string{"create"},
	Example: `gscloud ip add -6`,
	Short:   "Create a new IP address",
	Long: `Create a new IP address object.

Examples:

Create a new IPv6 address with a PTR entry and a name:

gscloud ip add -6 --name test --reverse-dns=example.com

Create a new IPv4 address:

gscloud ip add -4

`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if v4 && v6 {
			log.Fatal("No family given. Use either -4 or -6")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		family := gsclient.IPv4Type
		if v6 {
			family = gsclient.IPv6Type
		}
		ipOp := rt.IPOperator()
		ctx := context.Background()
		ipAddress, err := ipOp.CreateIP(ctx, gsclient.IPCreateRequest{
			Family:     family,
			Failover:   failover,
			Name:       name,
			ReverseDNS: reverseDNS,
		})
		if err != nil {
			log.Fatalf("Adding IPv%d address failed: %s", family, err)
		}
		fmt.Println("IP added:", ipAddress.IP)
	},
}

func init() {
	ipLsCmd.PersistentFlags().BoolVarP(&v4, "v4", "4", false, "IPv4 only")
	ipLsCmd.PersistentFlags().BoolVarP(&v6, "v6", "6", false, "IPv6 only")

	ipAddCmd.PersistentFlags().BoolVarP(&v4, "v4", "4", false, "Add a new IPv4 address")
	ipAddCmd.PersistentFlags().BoolVarP(&v6, "v6", "6", false, "Add a new IPv6 address")
	ipAddCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Optional name of the IP address being created. Can be omitted")
	ipAddCmd.PersistentFlags().BoolVarP(&failover, "failover", "", false, "Enable failover. If given, IP is no longer available for DHCP and cannot be assigned")
	ipAddCmd.PersistentFlags().StringVarP(&reverseDNS, "reverse-dns", "", "", "Optional reverse DNS entry for the IP address")

	ipCmd.AddCommand(ipLsCmd, ipRmCmd, ipAddCmd)
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
