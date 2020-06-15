package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gscloud/render"
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
			panic(err)
		}
		var storage [][]string
		if !jsonFlag {
			heading := []string{"name", "capacity", "changetime", "status", "id"}
			for _, stor := range storages {
				fill := [][]string{
					{
						stor.Properties.Name,
						strconv.FormatInt(int64(stor.Properties.Capacity), 10),
						strconv.FormatInt(int64(stor.Properties.ChangeTime.Hour()), 10),
						stor.Properties.Status,
						stor.Properties.ObjectUUID,
					},
				}
				storage = append(storage, fill...)
			}
			if idFlag {
				rowsToDisplay = len(heading)
			}
			render.Table(out, heading[:rowsToDisplay], storage)

		} else {
			render.AsJSON(out, storages)
		}
		fmt.Print(out)
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
