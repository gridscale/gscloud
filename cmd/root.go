package cmd

import (
	"fmt"
	"os"

	"github.com/gridscale/gscloud/render"
	"github.com/gridscale/gscloud/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootCmdFlags struct {
	configFile string
	account    string
	json       bool
	quiet      bool
	debug      bool
}

var (
	rootFlags  rootCmdFlags
	renderOpts render.Options
	rt         *runtime.Runtime
)

const (
	defaultAPIURL = "https://api.gridscale.io"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gscloud _object_ _verb_",
	Short: "the CLI for the gridscale API",
	Long: fmt.Sprintf(`gscloud lets you manage objects in the gridscale API via the command line. It provides a command line comparable to Docker-CLI that allows you to create, manipulate, and remove objects on gridscale.io.

Commands are given usually in the form of 'gscloud object verb'. For example, to list all servers you would do 'gscloud server ls'. Likewise, to list all storages you would do 'gscloud storage ls'.

Output, if any, is usually in the form of a table. You can pass --json to print output as JSON if you wish to do so.

To configure access to your projects via the API a YAML configuration file is used. See gscloud-make-config(1) and --config for more.

# EXAMPLES

List all servers available:

	$ gscloud server ls

    ID                                    NAME    CORE  MEM  CHANGED                    POWER
    37d53278-8e5f-47e1-a63f-54513e4b4d53  test-1  1     1    2020-11-17T08:48:22+01:00  off
    b0dd8d71-8f8d-46c1-8985-ce4b6dc37ecc  test-2  1     1    2020-11-20T11:44:58+01:00  off

Power all servers on:

    $ gscloud server ls --quiet | while read s; do
        gscloud server on $s
	done

Get the list of storages as JSON:

	$ gscloud --json storage ls | jq

	[
	  {
		"storage": {
		  "change_time": "2020-11-17T06:32:30Z",
		  "location_iata": "fra",
		  "status": "active",
		  "license_product_no": 0,
		  "location_country": "de",
		  "usage_in_minutes": 59230,
		  "last_used_template": "8d1bb5dc-7c37-4c90-8529-d2aaac75d812",
		  "current_price": 99.99,
		  "capacity": 10,
		  "location_uuid": "45ed677b-3702-4b36-be2a-a2eab9827950",
		  "storage_type": "storage",
		  "parent_uuid": "",
		  "name": "test-1",
		  "location_name": "de/fra",
		  "object_uuid": "479b7973-376a-4b23-98fc-50e94131a6e3",
		  "snapshots": [],
		  "relations": {
			"servers": [
			  {
				"bootdevice": true,
				"target": 0,
				"controller": 0,
				"bus": 0,
				"object_uuid": "37d53278-8e5f-47e1-a63f-54513e4b4d53",
				"lun": 0,
				"create_time": "2020-11-17T06:32:30Z",
				"object_name": "test-1"
			  }
			],
			"snapshot_schedules": []
		  },
		  "labels": [],
		  "create_time": "2020-11-17T06:32:30Z"
		}
	  }
	]

# FILES

%s user API access configuration
`, runtime.ConfigPathWithoutUser()),
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
		cobra.OnInitialize(initConfig, initRuntime, initLogging)
	}

	rootCmd.PersistentFlags().StringVar(&rootFlags.configFile, "config", runtime.ConfigPathWithoutUser(), fmt.Sprintf("Path to configuration file"))
	rootCmd.PersistentFlags().StringVarP(&rootFlags.account, "account", "", "default", "Specify the account used")
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.json, "json", "j", false, "Print JSON to stdout instead of a table")
	rootCmd.PersistentFlags().BoolVarP(&renderOpts.NoHeader, "noheading", "", false, "Do not print column headings")
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.quiet, "quiet", "q", false, "Print only object IDs")
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.debug, "debug", "", false, "Debug mode")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootFlags.configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(runtime.ConfigPath())
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
	theRuntime, err := runtime.NewRuntime(*conf, rootFlags.account)
	if err != nil {
		log.Fatal(err)
	}
	rt = theRuntime
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	if rootFlags.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
