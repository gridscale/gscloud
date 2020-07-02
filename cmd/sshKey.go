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
	sshKeyMainAction = iota
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
	nameFlag, fileFlag bool
)

func produceSSHKeyCmdRunFunc(g sshKeysOperator, action int) cmdRunFunc {
	switch action {
	case sshKeyMainAction:
		return func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			out := new(bytes.Buffer)
			sshkeys, err := g.GetSshkeyList(ctx)
			if err != nil {
				log.Error("Couldn't get SSH-keys:", err)
				return
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
			if nameFlag && fileFlag {
				ctx := context.Background()
				publicKey, err := ioutil.ReadFile(args[1])
				if err != nil {
					log.Error("Failed to read public-key from "+args[1], err)
				}
				key, err := g.CreateSshkey(ctx, gsclient.SshkeyCreateRequest{
					Name:   args[0],
					Sshkey: string(publicKey),
				})
				if err != nil {
					log.Error("Create SSH-key has failed with error", err)
					return
				}
				log.WithFields(log.Fields{
					"sshkey_uuid": key.ObjectUUID,
				}).Infof("SSH-key [%s] successfully created", args[0])
			}
		}

	case sshKeyDeleteAction:
		return func(cmd *cobra.Command, args []string) {
			if nameFlag {
				ctx := context.Background()
				sshkeys, err := g.GetSshkeyList(ctx)
				if err != nil {
					log.Error("Couldn't get SSH-keys:", err)
					return
				}
				for _, key := range sshkeys {
					if args[0] == key.Properties.ObjectUUID || args[0] == key.Properties.Name {
						err := g.DeleteSshkey(ctx, key.Properties.ObjectUUID)
						if err != nil {
							log.Error("Delete SSH-key has failed with error", err)
							return
						}
						log.Infof("SSH-key [%s] successfully removed", args[0])
					}
				}
			}
		}

	default:
	}
	return nil
}

func initSSHKeyCmd() {
	var sshKeyCmd = &cobra.Command{
		Use:   "ssh-key",
		Short: "Print ssh-key list",
		Long:  `Print all ssh-key information`,
		Run:   produceSSHKeyCmdRunFunc(client, sshKeyMainAction),
	}
	var addCmd, removeCmd = &cobra.Command{
		Use:   "add",
		Short: "add ssh-key",
		Long:  `Add ssh-key via file`,
		Args:  cobra.MinimumNArgs(2),
		Run:   produceSSHKeyCmdRunFunc(client, sshKeyAddAction),
	}, &cobra.Command{
		Use:   "remove",
		Short: "remove ssh-key",
		Long:  `Remove ssh-key via name or id`,
		Args:  cobra.MinimumNArgs(1),
		Run:   produceSSHKeyCmdRunFunc(client, sshKeyDeleteAction),
	}
	sshKeyCmd.AddCommand(addCmd, removeCmd)
	sshKeyCmd.PersistentFlags().BoolVarP(&nameFlag, "name", "n", false, "Set ssh-key name")
	sshKeyCmd.PersistentFlags().BoolVarP(&fileFlag, "file", "f", false, "Read ssh-key from file")
	rootCmd.AddCommand(sshKeyCmd)
}
