package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Operations on storages",
	Long:  `List, create, or remove storages.`,
}

var storageLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List storages",
	Long:    `List storage objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		storageOp := rt.StorageOperator()
		ctx := context.Background()
		out := new(bytes.Buffer)
		storages, err := storageOp.GetStorageList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get storage list: %s", err)
		}
		var rows [][]string
		if !jsonFlag {
			heading := []string{"id", "name", "capacity", "changetime", "status"}
			for _, storage := range storages {
				fill := [][]string{
					{
						storage.Properties.ObjectUUID,
						storage.Properties.Name,
						strconv.FormatInt(int64(storage.Properties.Capacity), 10),
						storage.Properties.ChangeTime.Local().Format(time.RFC3339),
						storage.Properties.Status,
					},
				}
				rows = append(rows, fill...)
			}
			if quietFlag {
				for _, info := range rows {
					fmt.Println(info[0])
				}
				return
			}
			render.AsTable(out, heading, rows, renderOpts)
		} else {
			for _, storage := range storages {
				render.AsJSON(out, storage)
			}
		}
		fmt.Print(out)
	},
}

var storageRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove storage",
	Long:    `Remove an existing storage.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		storageOp := rt.StorageOperator()
		ctx := context.Background()
		err := storageOp.DeleteStorage(ctx, args[0])
		if err != nil {
			log.Fatalf("Removing storage failed: %s", err)
		}
	},
}

func init() {
	storageCmd.AddCommand(storageLsCmd, storageRmCmd)
	rootCmd.AddCommand(storageCmd)
}
