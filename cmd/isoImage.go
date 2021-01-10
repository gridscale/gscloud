package cmd

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type isoImageCmdFlags struct {
	name string
}

var (
	isoImageFlags isoImageCmdFlags
)

var isoImageCmd = &cobra.Command{
	Use:   "iso-image",
	Short: "Operations on ISO images",
	Long:  `List, create, or remove ISO images.`,
}

var isoImageLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List images",
	Long:    `List ISO image objects.`,
	Run: func(cmd *cobra.Command, args []string) {

		imageOp := rt.ISOImageOperator()
		ctx := context.Background()
		out := new(bytes.Buffer)
		images, err := imageOp.GetISOImageList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get list of ISO images: %s", err)
		}
		var rows [][]string
		if !rootFlags.json {
			heading := []string{"id", "name", "changed", "private", "source url"}
			for _, image := range images {
				var private string
				if image.Properties.Private {
					private = "yes"
				} else {
					private = "no"
				}
				fill := [][]string{
					{
						image.Properties.ObjectUUID,
						image.Properties.Name,
						image.Properties.ChangeTime.Local().Format(time.RFC3339),
						private,
						image.Properties.SourceURL,
					},
				}
				rows = append(rows, fill...)
			}
			if rootFlags.quiet {
				for _, info := range rows {
					fmt.Println(info[0])
				}
			}
			render.AsTable(out, heading, rows, renderOpts)
		} else {
			render.AsJSON(out, images)
		}
		fmt.Print(out)
	},
}

func init() {
	isoImageCmd.AddCommand(isoImageLsCmd)
	rootCmd.AddCommand(isoImageCmd)
}
