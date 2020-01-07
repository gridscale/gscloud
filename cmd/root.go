package cmd

import (
	"fmt"
	"os"

	"github.com/kirsle/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	account string
	client  *gsclient
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gscloud",
	Short: "is the command line interface for the gridscale cloud.",
	Long:  "gscloud is the command line interface for the gridscale cloud.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("configuration file, default %s/config.yaml", configPath()))
	rootCmd.PersistentFlags().StringVar(&account, "account", "", "the account used, 'default' if none given")

	rootCmd.AddCommand(kubernetesCmd)
	kubernetesCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(execCredentialCmd)

}

// configPath construct platform specific path to the configuration file.
// - on Linux: $XDG_CONFIG_HOME or $HOME/.config
// - on macOS: $HOME/Library/Application Support
// - on Windows: %APPDATA% or "C:\\Users\\%USER%\\AppData\\Roaming"
func configPath() string {
	return configdir.LocalConfig("gridscale")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default paths.
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configPath())
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Not found. Disregard
		} else if _, ok := err.(*os.PathError); ok && contains(os.Args, "make-config") {
			// --config given along with make-config → we're about to create that file. Disregard
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	
	if account == "" {
		account = "default"
	}
	client = newCliClient(account)
	if client == nil {
		os.Exit(1)
	}
}

// contains tests whether string e is in slice s.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
