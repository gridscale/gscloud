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

var (
	nameFlag, fileFlag string
)

var sshKeyCmd = &cobra.Command{
	Use:   "ssh-key",
	Short: "Operations on SSH keys",
	Long:  `List, create, or remove SSH keys.`,
}

var sshKeyLsCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List SSH keys",
	Long:    `List SSH key objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		out := new(bytes.Buffer)
		op := rt.SSHKeyOperator()
		sshkeys, err := op.GetSshkeyList(ctx)
		if err != nil {
			log.Fatalf("Couldn't get SSH key list: %s", err)
		}
		var rows [][]string
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
				rows = append(rows, fill...)
			}

			if quietFlag {
				for _, info := range rows {
					fmt.Println(info[0])
				}
				return
			}
			render.AsTable(out, heading, rows, renderOpts)
		} else {
			render.AsJSON(out, sshkeys)
		}
		fmt.Print(out)
	},
}

var sshKeyAddCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "Add a new SSH key",
	Long:  `Create a new SSH key.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		op := rt.SSHKeyOperator()
		_, err = op.CreateSshkey(ctx, gsclient.SshkeyCreateRequest{
			Name:   nameFlag,
			Sshkey: string(publicKey),
		})
		if err != nil {
			log.Fatalf("Creating SSH key failed: %s", err)
		}
	},
}

var sshKeyRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove SSH key",
	Long:    `Remove an existing SSH key.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		op := rt.SSHKeyOperator()
		err := op.DeleteSshkey(ctx, args[0])
		if err != nil {
			log.Fatalf("Removing SSH key failed: %s", err)
		}
	},
}

func init() {
	sshKeyAddCmd.PersistentFlags().StringVarP(&nameFlag, "name", "n", "", "Name of the new key")
	sshKeyAddCmd.PersistentFlags().StringVarP(&fileFlag, "file", "f", "", "Path to public key file")

	sshKeyCmd.AddCommand(sshKeyLsCmd, sshKeyAddCmd, sshKeyRmCmd)
	rootCmd.AddCommand(sshKeyCmd)
}
