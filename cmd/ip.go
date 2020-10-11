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
	v4Only bool
	v6Only bool
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
		if v4Only && v6Only {
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
				if v4Only && addr.Properties.Family == 6 {
					continue
				}

				if v6Only && addr.Properties.Family == 4 {
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

func init() {
	ipLsCmd.PersistentFlags().BoolVarP(&v4Only, "v4-only", "4", false, "IPv4 family only")
	ipLsCmd.PersistentFlags().BoolVarP(&v6Only, "v6-only", "6", false, "IPv6 family only")

	ipCmd.AddCommand(ipLsCmd, ipRmCmd)
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
