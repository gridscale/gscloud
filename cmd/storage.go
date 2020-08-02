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
	Short: "Operations on storages",
	Long:  `List, create, or remove storages.`,
}

var storageListCmd = &cobra.Command{
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
			for _, storage := range storages {
				render.AsJSON(out, storage)
			}
		}
		fmt.Print(out)
	},
}

var storageRemoveCmd = &cobra.Command{
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
			log.Fatalf("Removing Storage failed: %s", err)
		}
	},
}

func init() {
	storageCmd.AddCommand(storageListCmd, storageRemoveCmd)
	rootCmd.AddCommand(storageCmd)
}
