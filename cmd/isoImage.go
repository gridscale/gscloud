package cmd

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type isoImageCmdFlags struct {
	name      string
	sourceURL string
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

var isoImageRmCmd = &cobra.Command{
	Use:     "rm [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove ISO image",
	Long:    `Remove an existing ISO image.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageOp := rt.ISOImageOperator()
		ctx := context.Background()
		err := imageOp.DeleteISOImage(ctx, args[0])
		if err != nil {
			log.Fatalf("Removing ISO image failed: %s", err)
		}
	},
}

var isoImageCreateCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: "Create a private ISO image",
	Long: `Create a new private ISO image.

# EXAMPLES

Create a Fedora CoreOS image:

	$ gscloud iso-image create \
	   --name="Fedora CoreOS" \
	   --source-url="https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/33.20201214.3.1/x86_64/fedora-coreos-33.20201214.3.1-live.x86_64.iso"
`,
	Run: func(cmd *cobra.Command, args []string) {
		imageOp := rt.ISOImageOperator()
		ctx := context.Background()
		image, err := imageOp.CreateISOImage(ctx, gsclient.ISOImageCreateRequest{
			Name:      isoImageFlags.name,
			SourceURL: isoImageFlags.sourceURL,
		})
		if err != nil {
			log.Fatalf("Creating ISO image failed: %s", err)
		}
		fmt.Println("Image created:", image.ObjectUUID)
	},
}

func init() {
	isoImageCreateCmd.Flags().StringVar(&isoImageFlags.name, "name", "", "Name of the image")
	isoImageCreateCmd.MarkFlagRequired("name")
	isoImageCreateCmd.Flags().StringVar(&isoImageFlags.sourceURL, "source-url", "", "URL from where the image is downloaded")
	isoImageCreateCmd.MarkFlagRequired("source-url")

	isoImageCmd.AddCommand(isoImageLsCmd, isoImageRmCmd, isoImageCreateCmd)
	rootCmd.AddCommand(isoImageCmd)
}
