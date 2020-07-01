package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Print storage list",
	Long:  `Print all storage information`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		storages, err := client.GetStorageList(ctx)
		if err != nil {
			log.Error("Couldn't get Storageinfo", err)
			return
		}
		var storageinfo [][]string
		if !jsonFlag {
			heading := []string{"id", "name", "capacity", "changetime", "status"}
			for _, stor := range storages {
				fill := [][]string{
					{
						stor.Properties.ObjectUUID,
						stor.Properties.Name,
						strconv.FormatInt(int64(stor.Properties.Capacity), 10),
						strconv.FormatInt(int64(stor.Properties.ChangeTime.Hour()), 10),
						stor.Properties.Status,
					},
				}
				storageinfo = append(storageinfo, fill...)
			}
			if quietFlag {
				for _, info := range storageinfo {
					fmt.Println(info[0])
				}
				return
			}
			render.Table(out, heading[:], storageinfo)

		} else {
			render.AsJSON(out, storages)
		}
		fmt.Print(out)
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
