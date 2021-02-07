package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/render"
	"github.com/spf13/cobra"
)

type sshKeyCmdFlags struct {
	name       string
	pubKeyFile string
}

var (
	sshKeyFlags sshKeyCmdFlags
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		out := new(bytes.Buffer)
		op := rt.SSHKeyOperator()
		sshkeys, err := op.GetSshkeyList(ctx)
		if err != nil {
			return NewError(cmd, "Could not get SSH key list", err)
		}
		var rows [][]string
		if !rootFlags.json {
			heading := []string{"id", "name", "key", "user", "created"}
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

			if rootFlags.quiet {
				for _, info := range rows {
					fmt.Println(info[0])
				}
				return nil
			}
			render.AsTable(out, heading, rows, renderOpts)
		} else {
			render.AsJSON(out, sshkeys)
		}
		fmt.Print(out)
		return nil
	},
}

var sshKeyAddCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "Add a new SSH key",
	Long:  `Create a new SSH key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		publicKey, err := ioutil.ReadFile(sshKeyFlags.pubKeyFile)
		if err != nil {
			return NewError(cmd, "Error reading file", err)
		}
		ctx := context.Background()
		op := rt.SSHKeyOperator()
		_, err = op.CreateSshkey(ctx, gsclient.SshkeyCreateRequest{
			Name:   sshKeyFlags.name,
			Sshkey: string(publicKey),
		})
		if err != nil {
			return NewError(cmd, "Creating SSH key failed", err)
		}
		return nil
	},
}

var sshKeyRmCmd = &cobra.Command{
	Use:     "rm [flags] [ID]",
	Aliases: []string{"remove"},
	Short:   "Remove SSH key",
	Long:    `Remove an existing SSH key.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		op := rt.SSHKeyOperator()
		err := op.DeleteSshkey(ctx, args[0])
		if err != nil {
			return NewError(cmd, "Removing SSH key failed", err)
		}
		return nil
	},
}

func init() {
	sshKeyAddCmd.PersistentFlags().StringVarP(&sshKeyFlags.name, "name", "n", "", "Name of the new key")
	sshKeyAddCmd.MarkFlagRequired("name")
	sshKeyAddCmd.PersistentFlags().StringVarP(&sshKeyFlags.pubKeyFile, "file", "f", "", "Path to public key file")
	sshKeyAddCmd.MarkFlagRequired("file")

	sshKeyCmd.AddCommand(sshKeyLsCmd, sshKeyAddCmd, sshKeyRmCmd)
	rootCmd.AddCommand(sshKeyCmd)
}
