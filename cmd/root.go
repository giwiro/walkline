package cmd

import (
    "fmt"
    "github.com/giwiro/walkline/utils"
    "github.com/spf13/viper"
    "log"
    "os"

    "github.com/giwiro/walkline/version"
    "github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:     "walkline",
    Version: version.Version,
    CompletionOptions: cobra.CompletionOptions{
        DisableDefaultCmd: true,
    },
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        var verbose = utils.GetFlagBooleanValue(cmd, "verbose", false)

        if verbose == true {
            var configFile = viper.ConfigFileUsed()
            if configFile == "" {
                fmt.Println("Config file not found")
            } else {
                fmt.Println("Using config file:", viper.ConfigFileUsed())
            }
        }
    },
    Short: "Simplistic sql database migration tool",
    Long: `
               _ _    _            
              | | |  | (_)           
__      __ _ _| | | _| |_ _ __   ___ 
\ \ /\ / / _` + "` " + `| | |/ | | | '_ \ / _ \
 \ V  V | (_| | |   <| | | | | |  __/
  \_/\_/ \__,_|_|_|\_|_|_|_| |_|\___|
	Simplistic sql database migration tool
`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    // Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    log.SetFlags(0)

    rootCmd.PersistentFlags().StringP("url", "u", "", "sql database connection url")
    rootCmd.PersistentFlags().StringP("path", "p", "", "path of the migration files")
    rootCmd.PersistentFlags().StringP("schema", "s", "", "select the version table schema")
    rootCmd.PersistentFlags().BoolP("verbose", "v", false, "add verbosity")
}

func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory
        home, err := os.UserHomeDir()

        if err != nil {
            return
        }

        // Get working directory
        workingDir, err := utils.GetWorkingDir()

        if err != nil {
            return
        }

        if err != nil {
            log.Fatal(err)
        }

        // Search config in home directory with name ".cobra" (without extension).
        viper.AddConfigPath(workingDir)
        viper.AddConfigPath(home)
        viper.SetConfigName("walkline")
        viper.SetConfigType("yaml")

        _ = viper.ReadInConfig()
    }
}
