package cmd

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/utils"
	"github.com/spf13/cobra"
)

func unique(strs []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strs {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

var postgresCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "Operate managed PostgreSQL database",
	Long:  "Create, manipulate, and remove managed PostgreSQL database systems.",
}

var postgresReleasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "Returns the available PostgreSQL releases",
	Long:  "Returns the available PostgreSQL releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		out := new(bytes.Buffer)
		op := rt.PaaSOperator()
		paasTemplates, err := op.GetPaaSTemplateList(ctx)
		if err != nil {
			return NewError(cmd, "Could not get get list of PostgreSQL releases", err)
		}

		var releases []string
		for _, template := range paasTemplates {
			if template.Properties.Flavour == "postgres" {
				releases = append(releases, template.Properties.Release)
			}
		}
		releases = unique(releases)
		sort.Sort(sort.Reverse(utils.StringSorter(releases)))
		if !rootFlags.json {
			heading := []string{"releases"}
			var rows [][]string
			for _, rel := range releases {
				rows = append(rows, []string{rel})
			}
			render.AsTable(out, heading, rows, renderOpts)
			if rootFlags.quiet {
				for _, rel := range releases {
					fmt.Println(rel)
				}
				return nil
			}

		} else {
			render.AsJSON(out, releases)
		}
		fmt.Print(out)
		return nil
	},
}

func init() {
	postgresCmd.AddCommand(postgresReleasesCmd)
	rootCmd.AddCommand(postgresCmd)
}
