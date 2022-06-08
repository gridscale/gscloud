package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print project summary",
	Long: `Print information about the current project.

# EXAMPLES

Show summary for a given project:

	$ gscloud --project=dev@example.com info
	SETTING    VALUE
	Project    dev@example.com
	User ID    7ff8003b-55c5-45c5-bf0c-3746735a4f99
	API token  <redacted>
	URL        https://api.gridscale.io
	Getting information about used resources…
	OBJECT             COUNT
	Platform services  0
	Servers            18
	Storages           24
	IP addresses       2
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		type output struct {
			runtime.ProjectEntry
			ServerAgg  map[string]interface{} `json:"server"`
			StorageAgg map[string]interface{} `json:"storage"`
			IPAddrAgg  map[string]interface{} `json:"ip_address"`
			PaasAgg    map[string]interface{} `json:"platform_service"`
		}

		type objectCount struct {
			Obj string
			Agg map[string]interface{}
			Err error
		}

		conf := rt.Config()
		for _, project := range conf.Projects {
			projectName := rt.Project()
			if project.Name == projectName {
				if !rootFlags.json {
					out := new(bytes.Buffer)
					heading := []string{"setting", "value"}
					fill := [][]string{
						{"Project", project.Name},
						{"User ID", project.UserID},
						{"API token", project.Token},
						{"URL", project.URL},
					}
					var rows [][]string
					rows = append(rows, fill...)
					render.AsTable(out, heading, rows, renderOpts)
					fmt.Print(out)
				}

				fmt.Fprintln(os.Stderr, "Getting information about used resources…")
				client := rt.Client()

				funcs := map[string]func(context.Context, *gsclient.Client) (map[string]interface{}, error){
					"Servers": func(ctx context.Context, c *gsclient.Client) (map[string]interface{}, error) {
						objs, err := c.GetServerList(ctx)
						if err != nil {
							return nil, err
						}
						mem := 0
						cores := 0
						for _, obj := range objs {
							mem += obj.Properties.Memory
							cores += obj.Properties.Cores
						}
						return map[string]interface{}{
								"count":  len(objs),
								"memory": mem,
								"cores":  cores,
							},
							nil
					},
					"Storages": func(ctx context.Context, c *gsclient.Client) (map[string]interface{}, error) {
						objs, err := c.GetStorageList(ctx)
						if err != nil {
							return nil, err
						}
						capacity := 0
						for _, obj := range objs {
							capacity += obj.Properties.Capacity
						}
						return map[string]interface{}{
								"count":    len(objs),
								"capacity": capacity,
							},
							nil
					},
					"IP addresses": func(ctx context.Context, c *gsclient.Client) (map[string]interface{}, error) {
						objs, err := c.GetIPList(ctx)
						if err != nil {
							return nil, err
						}
						return map[string]interface{}{
								"count": len(objs),
							},
							nil
					},
					"Platform services": func(ctx context.Context, c *gsclient.Client) (map[string]interface{}, error) {
						objs, err := c.GetPaaSServiceList(ctx)
						if err != nil {
							return nil, err
						}
						return map[string]interface{}{
								"count": len(objs),
							},
							nil
					},
				}

				var wg sync.WaitGroup
				ch := make(chan objectCount)

				for k, v := range funcs {
					wg.Add(1)
					cCopy := context.Background()
					go func(obj string, f func(context.Context, *gsclient.Client) (map[string]interface{}, error)) {
						defer wg.Done()

						agg, err := f(cCopy, client)
						if err != nil {
							ch <- objectCount{obj, nil, NewError(cmd, fmt.Sprintf("Could not get %s", obj), err)}
						}
						ch <- objectCount{obj, agg, nil}
					}(k, v)
				}

				go func() {
					wg.Wait()
					close(ch)
				}()

				out := new(bytes.Buffer)
				if !rootFlags.json {
					heading := []string{"object", "count"}
					var rows [][]string
					for v := range ch {
						if v.Err != nil {
							return v.Err
						}
						count := v.Agg["count"].(int)
						rows = append(rows, []string{v.Obj, strconv.Itoa(count)})
					}
					render.AsTable(out, heading, rows, renderOpts)
				} else {
					m := map[string]map[string]interface{}{}
					for v := range ch {
						if v.Err != nil {
							return v.Err
						}
						m[v.Obj] = v.Agg
					}

					jsonOutput := output{
						ProjectEntry: project,
						ServerAgg:    m["Servers"],
						StorageAgg:   m["Storages"],
						IPAddrAgg:    m["IP addresses"],
						PaasAgg:      m["Platform services"],
					}
					render.AsJSON(out, jsonOutput)
				}
				fmt.Print(out)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Hide some global persistent flags
	origHelpFunc := infoCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "info" || (cmd.HasParent() && cmd.Parent().Name() == "info") {
			cmd.Flags().MarkHidden("config")
			cmd.Flags().MarkHidden("quiet")
			cmd.Flags().MarkHidden("noheading")
		}
		origHelpFunc(cmd, args)
	})
}
