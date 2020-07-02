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

// istorageRetriever is used for testing purpose,
// we can mock data return from the gsclient via interface.
type storageGetter interface {
	GetStorageList(ctx context.Context) ([]gsclient.Storage, error)
}

// produceStorageCmdRunFunc takes an instance of a struct that implements `storageGetter`
// returns a `cmdRunFun`
func produceStorageCmdRunFunc(g storageGetter) cmdRunFunc {
	return func(cmd *cobra.Command, args []string) {
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
	}
}

func init() {
	var storageCmd = &cobra.Command{
		Use:   "storage",
		Short: "Print storage list",
		Long:  `Print all storage information`,
		Run:   produceStorageCmdRunFunc(client),
	}
	rootCmd.AddCommand(storageCmd)
}
