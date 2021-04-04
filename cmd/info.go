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
	Short: "Print account summary",
	Long:  `Print information about the current accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		type output struct {
			runtime.AccountEntry
			ServerCount  int `json:"server_count"`
			StorageCount int `json:"storage_count"`
			IPAddrCount  int `json:"ip_address_count"`
			PaaSCount    int `json:"platform_service_count"`
		}

		type objectCount struct {
			Obj   string
			Count int
			Err   error
		}

		conf := rt.Config()
		for _, account := range conf.Accounts {
			accountName := rt.Account()
			if account.Name == accountName {
				if !rootFlags.json {
					out := new(bytes.Buffer)
					heading := []string{"setting", "value"}
					fill := [][]string{
						{"Account", account.Name},
						{"User ID", account.UserID},
						{"API token", account.Token},
						{"URL", account.URL},
					}
					var rows [][]string
					rows = append(rows, fill...)
					render.AsTable(out, heading, rows, renderOpts)
					fmt.Print(out)
				}

				fmt.Fprintf(os.Stderr, "Getting information about used resources…\n")
				client := rt.Client()

				funcs := map[string]func(context.Context, *gsclient.Client) (int, error){
					"Servers": func(ctx context.Context, c *gsclient.Client) (int, error) {
						objs, err := c.GetServerList(ctx)
						return len(objs), err
					},
					"Storages": func(ctx context.Context, c *gsclient.Client) (int, error) {
						objs, err := c.GetStorageList(ctx)
						return len(objs), err
					},
					"IP addresses": func(ctx context.Context, c *gsclient.Client) (int, error) {
						objs, err := c.GetIPList(ctx)
						return len(objs), err
					},
					"Platform services": func(ctx context.Context, c *gsclient.Client) (int, error) {
						objs, err := c.GetPaaSServiceList(ctx)
						return len(objs), err
					},
				}

				var wg sync.WaitGroup
				ch := make(chan objectCount)

				go func() {
					wg.Wait()
					close(ch)
				}()

				for k, v := range funcs {
					wg.Add(1)
					cCopy := context.Background()
					go func(obj string, f func(context.Context, *gsclient.Client) (int, error)) {
						defer wg.Done()

						count, err := f(cCopy, client)
						if err != nil {
							ch <- objectCount{obj, 0, NewError(cmd, fmt.Sprintf("Could not get %s", obj), err)}
						}
						ch <- objectCount{obj, count, nil}
					}(k, v)
				}

				out := new(bytes.Buffer)
				if !rootFlags.json {
					heading := []string{"object", "count"}
					var rows [][]string
					for v := range ch {
						if v.Err != nil {
							return v.Err
						} else {
							rows = append(rows, []string{v.Obj, strconv.Itoa(v.Count)})
						}
					}
					render.AsTable(out, heading, rows, renderOpts)
				} else {
					m := map[string]int{}
					for v := range ch {
						if v.Err != nil {
							return v.Err
						} else {
							m[v.Obj] = v.Count
						}
					}
					jsonOutput := output{
						AccountEntry: account,
						ServerCount:  m["Servers"],
						StorageCount: m["Storages"],
						IPAddrCount:  m["IP addresses"],
						PaaSCount:    m["Platform services"],
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
