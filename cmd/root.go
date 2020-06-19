package cmd

import (
	"fmt"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	account       string
	client        *gsclient.Client
	jsonFlag      bool
	idFlag        bool
	rowsToDisplay = 4
)

const (
	requestBase                    = "/requests/"
	apiPaasServiceBase             = "/objects/paas/services"
	defaultAPIURL                  = "https://api.gridscale.io"
	bodyType                       = "application/json"
	requestDoneStatus              = "done"
	requestFailStatus              = "failed"
	defaultCheckRequestTimeoutSecs = 120
	defaultDelayIntervalMilliSecs  = 500
	requestUUIDHeaderParam         = "X-Request-Id"
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
	cobra.OnInitialize(initConfig, initClient)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("configuration file, default %s", cliConfigPath()))
	rootCmd.PersistentFlags().StringVar(&account, "account", "", "the account used, 'default' if none given")
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Print JSON to stdout instead of a table")
	rootCmd.PersistentFlags().BoolVarP(&idFlag, "id", "i", false, "Include ID column")

	rootCmd.AddCommand(kubernetesCmd)
	kubernetesCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(execCredentialCmd)
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
		viper.AddConfigPath(cliConfigPath())
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Not found. Disregard
		} else if _, ok := err.(*os.PathError); ok && commandWithoutConfig(os.Args) {
			// --config given along with make-config → we're about to create that file. Disregard
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// initClient initializes the client for a given account.
func initClient() {
	if account == "" {
		account = "default"
	}

	client = newCliClient(account)
	if client == nil {
		os.Exit(1)
	}

}

// commandWithoutConfig return true if current command does not need a config file.
// Called from within a cobra initializer function. Unfortunately there is no
// way of getting the current command from an cobra initializer so we scan the
// command line again.
func commandWithoutConfig(cmdLine []string) bool {
	var noConfigNeeded = []string{
		"make-config", "version",
	}
	for _, cmd := range noConfigNeeded {
		if contains(cmdLine, cmd) {
			return true
		}
	}
	return false
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
