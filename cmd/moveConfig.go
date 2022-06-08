package cmd

import (
	"fmt"

	"github.com/gridscale/gscloud/runtime"
	"github.com/gridscale/gscloud/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type moveConfigCmdFlags struct {
	force bool
}

var moveConfigFlags moveConfigCmdFlags

var moveConfigCmd = &cobra.Command{
	Use:   "move-config",
	Short: "Move an old config file to the current path",
	Long: fmt.Sprintf(`Move an old config file (like at %s) to the current position (%s), 
 while updating it's contents to match the current format. Doesn't override the old config when --force is not given`,
		runtime.OldConfigPathWithoutUser(), runtime.ConfigPathWithoutUser()),

	RunE: func(cmd *cobra.Command, args []string) error {

		viper.SetConfigFile(runtime.OldConfigPath())

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Infoln("Old config not found, stopping")
				return nil
			}
			return err
		}

		if utils.FileExists(runtime.ConfigPath()) && !moveConfigFlags.force {
			log.Errorf("This action would overwrite %s. Run with --force to force this\n", runtime.ConfigPath())
			return nil
		}

		conf, err := runtime.ParseConfig()

		if err != nil {
			return err
		}

		return runtime.WriteConfig(conf, runtime.ConfigPath())
	},
}

func init() {
	moveConfigCmd.Flags().BoolVarP(&moveConfigFlags.force, "force", "f", false, "Force overwriting old config files")

	rootCmd.AddCommand(moveConfigCmd)
}
