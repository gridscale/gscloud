package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Operations on templates",
	Long:  `List templates.`,
}

var templateLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List templates",
	Long:    `List template objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateOp := rt.TemplateOperator()
		ctx := context.Background()
		out := new(bytes.Buffer)
		templates, err := templateOp.GetTemplateList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get template list: %s", err)
		}
		var rows [][]string
		if jsonFlag {
			render.AsJSON(out, templates)
		} else {
			heading := []string{"id", "name", "capacity", "changed", "description"}
			for _, template := range templates {
				fill := [][]string{
					{
						template.Properties.ObjectUUID,
						template.Properties.Name,
						strconv.FormatInt(int64(template.Properties.Capacity), 10),
						template.Properties.ChangeTime.Local().Format(time.RFC3339),
						template.Properties.Description,
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
		}
		fmt.Print(out)
	},
}

var templateRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove templates",
	Long:    `Remove a template by ID.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		storageOp := rt.TemplateOperator()
		ctx := context.Background()
		err := storageOp.DeleteTemplate(ctx, args[0])
		if err != nil {
			log.Fatalf("Removing template failed: %s", err)
		}
	},
}

func init() {
	templateCmd.AddCommand(templateLsCmd, templateRmCmd)
	rootCmd.AddCommand(templateCmd)
}
