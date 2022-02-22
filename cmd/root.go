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
	log.SetFlags(1)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.walkline.yaml)")
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .walkline.yaml)")
	rootCmd.PersistentFlags().StringP("flavor", "f", "", "sql database brand [postgresql]")
	rootCmd.PersistentFlags().StringP("url", "u", "", "sql database connection url")
	rootCmd.PersistentFlags().StringP("path", "p", "", "path of the migration files")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "add verbosity")

	// Bind with Viper
	/* err := viper.BindPFlag("flavor", rootCmd.PersistentFlags().Lookup("flavor"))
	if err != nil {
		fmt.Println(err)
	}
	err = viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	if err != nil {
		fmt.Println(err)
	} */
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


		if err := viper.ReadInConfig(); err == nil {
			if viper.GetBool("verbose") == true {
				fmt.Println("Using config file:", viper.ConfigFileUsed())
			}
		}
	}
}
