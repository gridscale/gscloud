package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type storageCmdFlags struct {
	name     string
	capacity int
	force    bool
}

var (
	storageFlags storageCmdFlags
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
	RunE: func(cmd *cobra.Command, args []string) error {
		storageOp := rt.StorageOperator()
		ctx := context.Background()
		out := new(bytes.Buffer)
		storages, err := storageOp.GetStorageList(ctx)
		if err != nil {
			return NewError(cmd, "Could not get storage list", err)
		}
		var rows [][]string
		if !rootFlags.json {
			heading := []string{"id", "name", "capacity", "changed", "status"}
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
			if rootFlags.quiet {
				for _, info := range rows {
					fmt.Println(info[0])
				}
				return nil
			}
			render.AsTable(out, heading, rows, renderOpts)
		} else {
			render.AsJSON(out, storages)
		}
		fmt.Print(out)
		return nil
	},
}

var storageSetCmd = &cobra.Command{
	Use:     "set [flags] ID",
	Example: `gscloud storage set --capacity 20 b3ec341c-1732-45b3-bc45-9a7fcebb363e`,
	Short:   "Update storage properties",
	Long: `Update properties of a storage object.

# EXAMPLES

Rename a storage object:

    $ gscloud storage set --name test-1 b3ec341c-1732-45b3-bc45-9a7fcebb363e

Shrink a storage:

    $ gscloud storage set --capacity 9 --force b3ec341c-1732-45b3-bc45-9a7fcebb363e
`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Name == "capacity" && f.Changed {
				if storageFlags.capacity < 1 {
					err = errors.New("expected storage capacity ≥ 1 GB")
				}
			}
		})
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		storageOp := rt.StorageOperator()
		updateReq := gsclient.StorageUpdateRequest{}
		if len(storageFlags.name) > 0 {
			updateReq.Name = storageFlags.name
		}
		if storageFlags.capacity > 0 {
			storage, err := storageOp.GetStorage(ctx, args[0])
			if err != nil {
				return NewError(cmd, "Could not set new capacity", err)
			}
			currentSize := storage.Properties.Capacity
			if storageFlags.capacity < currentSize {
				if !storageFlags.force {
					log.Printf("Downsizing can destroy your data. Re-run with --force to reduce storage size from %d GB to %d GB\n", currentSize, storageFlags.capacity)
					return nil
				}
			}
			updateReq.Capacity = storageFlags.capacity
		}
		err := storageOp.UpdateStorage(
			ctx,
			args[0],
			updateReq)
		if err != nil {
			return NewError(cmd, "Could not set property", err)
		}
		return nil
	},
}

var storageRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove storage",
	Long:    `Remove an existing storage.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storageOp := rt.StorageOperator()
		ctx := context.Background()
		err := storageOp.DeleteStorage(ctx, args[0])
		if err != nil {
			return NewError(cmd, "Deleting storage failed", err)
		}
		fmt.Fprintf(os.Stderr, "Removed %s\n", args[0])
		return nil
	},
}

func init() {
	storageSetCmd.PersistentFlags().StringVarP(&storageFlags.name, "name", "n", "", "Change name")
	storageSetCmd.PersistentFlags().IntVar(&storageFlags.capacity, "capacity", 0, "Change size (GB)")
	storageSetCmd.PersistentFlags().BoolVarP(&storageFlags.force, "force", "", false, "Force a potential destructive operation")

	storageCmd.AddCommand(storageLsCmd, storageSetCmd, storageRmCmd)
	rootCmd.AddCommand(storageCmd)
}
