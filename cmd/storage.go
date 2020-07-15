package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// storageOperator is used for testing purpose,
// we can mock data return from the gsclient via interface.
type storageOperator interface {
	GetStorageList(ctx context.Context) ([]gsclient.Storage, error)
	DeleteStorage(ctx context.Context, id string) error
	//	CreateStorage(ctx, capacity gsclient.StorageCreateRequest) error
}

// Storage action enums
const (
	storageListAction = iota
	storageDeleteAction
	storageCreateAction
)

// produceStorageCmdRunFunc takes an instance of a struct that implements `storageOperator`
// returns a `cmdRunFunc`
func produceStorageCmdRunFunc(o storageOperator, action int) cmdRunFunc {
	switch action {
	case storageListAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			out := new(bytes.Buffer)
			storages, err := o.GetStorageList(ctx)
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
		}
	case storageDeleteAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := o.DeleteStorage(ctx, args[0])
			if err != nil {
				log.Fatalf("Removing Storage failed: %s", err)
			}
		}
	case storageCreateAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			s := args[1]
			cap := 0
			if n, err := strconv.Atoi(s); err == nil {
				cap = n
			}
			cStorage, err := client.CreateStorage(ctx, gsclient.StorageCreateRequest{
				Name:        args[0],
				Capacity:    cap,
				StorageType: gsclient.InsaneStorageType,
			})
			if err != nil {
				log.Error("Create storage has failed with error", err)
				return
			}
			log.WithFields(log.Fields{
				"storage_uuid": cStorage.ObjectUUID,
			}).Info("Storage successfully created")
		}
	default:
	}
	return nil
}

// initStorageCmd adds storage cmd to the root cmd
func initStorageCmd() {
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
		Run:     produceStorageCmdRunFunc(client, storageListAction),
	}
	var removeCmd = &cobra.Command{
		Use:     "rm [flags] [ID]",
		Aliases: []string{"remove"},
		Short:   "Remove storage",
		Long:    `Remove an existing storage.`,
		Args:    cobra.ExactArgs(1),
		Run:     produceStorageCmdRunFunc(client, storageDeleteAction),
	}
	var createCmd = &cobra.Command{
		Use:   "create [Name] [Capacity]",
		Short: "Create storage",
		Long:  `Create a new storage.`,
		Args:  cobra.MaximumNArgs(3),
		Run:   produceStorageCmdRunFunc(client, storageCreateAction),
	}
	storageCmd.AddCommand(storageLsCmd, removeCmd, createCmd)
	rootCmd.AddCommand(storageCmd)
}
