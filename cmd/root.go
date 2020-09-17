package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	account    string
	rt         *runtime.Runtime
	jsonFlag   bool
	quietFlag  bool
	renderOpts render.Options
)

const (
	defaultAPIURL = "https://api.gridscale.io"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gscloud",
	Short: "the CLI for the gridscale cloud",
	Long: `gscloud lets you manage objects on gridscale.io via the command line. It
provides a command line comparable to Docker-CLI that allows you to create,
manipulate, and remove objects on gridscale.io.`,
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register following initializers only when we are not running tests.
	if !runtime.UnderTest() {
		cobra.OnInitialize(initConfig, initRuntime)
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", runtime.ConfigPath(), fmt.Sprintf("Specify a configuration file"))
	rootCmd.PersistentFlags().StringVarP(&account, "account", "", "default", "Specify the account used")
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Print JSON to stdout instead of a table")
	rootCmd.PersistentFlags().BoolVarP(&renderOpts.NoHeader, "noheading", "", false, "Do not print column headings")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Print only IDs of objects")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Use default paths.
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(runtime.ConfigPath())
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Not found. Disregard
		} else if _, ok := err.(*os.PathError); ok && runtime.CommandWithoutConfig(os.Args) {
			// --config given along with make-config → we're about to create that file. Disregard
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// initRuntime initializes the client for a given account.
func initRuntime() {
	conf, err := runtime.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	theRuntime, err := runtime.NewRuntime(*conf, account)
	if err != nil {
		log.Fatal(err)
	}
	rt = theRuntime
}
