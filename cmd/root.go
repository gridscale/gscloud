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

// Error is returned when something goes wrong in a command or sub-command.
type Error struct {
	Command *cobra.Command
	What    string
	Err     error
}

func (e *Error) Error() string { return e.What + ": " + e.Err.Error() }

// NewError constructs a new error.
func NewError(cmd *cobra.Command, what string, err error) *Error {
	return &Error{Command: cmd, What: what, Err: err}
}

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

# FILES

%s user API access configuration

# EXIT CODES

gscloud returns zero exit code on success, non-zero on failure. Following exit codes map to these failure modes:

    1. The requested command failed.
    2. Reading the configuration file failed.
    3. The configuration could not be parsed.
    4. The account specified does not exist in the configuration file.

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

# ENVIRONMENT

GRIDSCALE_ACCOUNT
	Specify the account used. Gets overriden by --account option

GRIDSCALE_UUID
	Specify the user id used. Overrides the value in the config file

GRIDSCALE_TOKEN
	Specify the API token used. Overrides the value in the config file

GRIDSCALE_URL
	Specify the URL of the API. Overrides the value in the config file

`, runtime.ConfigPathWithoutUser()),
	DisableAutoGenTag: true,
}

// Execute runs the subcommand. Execute adds all child commands to the root
// command and sets flags appropriately. In case of errors Execute prints the
// returned error to stderr and ends the process with a non-zero exit code.
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

	// Do not print usage or error strings in case of errors. Commands use RunE
	// and return errors. We print errors in Execute.
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	account, accountEnvPresent := os.LookupEnv("GRIDSCALE_ACCOUNT")
	if !accountEnvPresent {
		account = "default"
	}

	rootCmd.PersistentFlags().StringVar(&rootFlags.configFile, "config", runtime.ConfigPathWithoutUser(), "Path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.account, "account", "", account, "Specify the account used. Overrides GRIDSCALE_ACCOUNT environment variable")
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
		} else if _, ok := err.(*os.PathError); ok && CommandWithoutConfig(os.Args) {
			// --config given along with make-config → we're about to create that file. Disregard
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
}

// initRuntime initializes the client for a given account.
func initRuntime() {
	conf, err := runtime.ParseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	theRuntime, err := runtime.NewRuntime(*conf, rootFlags.account, CommandWithoutConfig(os.Args))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
	rt = theRuntime
}

type plainFormatter struct {
}

func (f *plainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func initLogging() {

	log.SetFormatter(&plainFormatter{})

	if rootFlags.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// CommandWithoutConfig return true if current command does not need a config file.
// Called from within a cobra initializer function. Unfortunately there is no
// way of getting the current command from an cobra initializer so we scan the
// command line again.
func CommandWithoutConfig(cmdLine []string) bool {
	noConfigNeeded := []string{
		"make-config", "version", "manpage", "completion",
	}

	foundCommand, _, _ := rootCmd.Find(cmdLine[1:])

	return contains(noConfigNeeded, foundCommand.Name())
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
