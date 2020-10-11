package cmd

import (
	"bytes"
	"context"
	"fmt"

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
		ipAdresses, err := ipOp.GetIPList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get list of IPs: %s", err)
		}
		var rows [][]string
		out := new(bytes.Buffer)
		if !jsonFlag {
			heading := []string{"IP", "assigned", "failover", "family", "reverse DNS", "ID"}
			for _, addr := range ipAdresses {
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
			render.AsJSON(out, ipAdresses)
		}
		fmt.Print(out)
	},
}

func init() {
	ipLsCmd.PersistentFlags().BoolVarP(&v4Only, "v4-only", "4", false, "IPv4 family only")
	ipLsCmd.PersistentFlags().BoolVarP(&v6Only, "v6-only", "6", false, "IPv6 family only")

	ipCmd.AddCommand(ipLsCmd)
	rootCmd.AddCommand(ipCmd)
}
