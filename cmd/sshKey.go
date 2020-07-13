package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SSH keys action enums
const (
	sshKeyListAction = iota
	sshKeyAddAction
	sshKeyDeleteAction
)

// sshKeysOperator is used for testing purpose,
// we can mock data return from the gsclient via interface.
type sshKeysOperator interface {
	GetSshkeyList(ctx context.Context) ([]gsclient.Sshkey, error)
	CreateSshkey(ctx context.Context, body gsclient.SshkeyCreateRequest) (gsclient.CreateResponse, error)
	DeleteSshkey(ctx context.Context, id string) error
}

var (
	nameFlag, fileFlag string
)

// produceSSHKeyCmdRunFunc takes an instance of a struct that implements `sshKeysOperator`
// returns a `cmdRunFunc`
func produceSSHKeyCmdRunFunc(o sshKeysOperator, action int) cmdRunFunc {
	switch action {
	case sshKeyListAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			out := new(bytes.Buffer)
			sshkeys, err := o.GetSshkeyList(ctx)
			if err != nil {
				log.Fatalf("Couldn't get SSH key list: %s", err)
			}
			var sshkeyinfo [][]string
			if !jsonFlag {
				heading := []string{"id", "name", "key", "user", "createtime"}
				for _, key := range sshkeys {
					fill := [][]string{
						{
							key.Properties.ObjectUUID,
							key.Properties.Name,
							key.Properties.Sshkey[:10] + "..." + key.Properties.Sshkey[len(key.Properties.Sshkey)-30:],
							key.Properties.UserUUID[:8],
							key.Properties.CreateTime.String()[:19],
						},
					}
					sshkeyinfo = append(sshkeyinfo, fill...)
				}

				if quietFlag {
					for _, info := range sshkeyinfo {
						fmt.Println(info[0])
					}
					return
				}
				render.Table(out, heading[:], sshkeyinfo)
			} else {
				render.AsJSON(out, sshkeys)
			}
			fmt.Print(out)
		}

	case sshKeyAddAction:
		return func(cmd *cobra.Command, args []string) {
			if !cmd.Flag("name").Changed {
				log.Fatalf("Mandatory flag missing: name")
			}
			if !cmd.Flag("file").Changed {
				log.Fatalf("Mandatory flag missing: file")
			}

			publicKey, err := ioutil.ReadFile(fileFlag)
			if err != nil {
				log.Fatalf("Error reading file: %s", err)
			}
			ctx := context.Background()
			_, err = o.CreateSshkey(ctx, gsclient.SshkeyCreateRequest{
				Name:   nameFlag,
				Sshkey: string(publicKey),
			})
			if err != nil {
				log.Fatalf("Creating SSH key failed: %s", err)
			}
		}

	case sshKeyDeleteAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := o.DeleteSshkey(ctx, args[0])
			if err != nil {
				log.Fatalf("Removing SSH key failed: %s", err)
			}
		}

	default:
	}
	return nil
}

func initSSHKeyCmd() {
	var sshKeyCmd = &cobra.Command{
		Use:   "ssh-key",
		Short: "Operations on SSH keys",
		Long:  `List, create, or remove SSH keys.`,
	}

	var lsCmd = &cobra.Command{
		Use:     "ls [flags]",
		Aliases: []string{"list"},
		Short:   "List SSH keys",
		Long:    `List SSH key objects.`,
		Run:     produceSSHKeyCmdRunFunc(client, sshKeyListAction),
	}

	var addCmd = &cobra.Command{
		Use:   "add [flags]",
		Short: "Add a new SSH key",
		Long:  `Create a new SSH key.`,
		Run:   produceSSHKeyCmdRunFunc(client, sshKeyAddAction),
	}
	addCmd.PersistentFlags().StringVarP(&nameFlag, "name", "n", "", "Name of the new key")
	addCmd.PersistentFlags().StringVarP(&fileFlag, "file", "f", "", "Path to public key file")

	var removeCmd = &cobra.Command{
		Use:     "rm [flags] [ID]",
		Aliases: []string{"remove"},
		Short:   "Remove SSH key",
		Long:    `Remove an existing SSH key.`,
		Args:    cobra.ExactArgs(1),
		Run:     produceSSHKeyCmdRunFunc(client, sshKeyDeleteAction),
	}

	sshKeyCmd.AddCommand(lsCmd, addCmd, removeCmd)
	rootCmd.AddCommand(sshKeyCmd)
}
