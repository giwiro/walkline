package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/giwiro/walkline/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Long: `Initializes the version table in the default schema`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("flags", cmd.Flags().Lookup("url"))
		err := cmd.MarkFlagRequired("url")
		if err != nil {
			return 
		}
		fmt.Println("pflags", cmd.PersistentFlags().Lookup("url"))
		err = viper.BindPFlag("url", cmd.Flags().Lookup("url"))
		if err != nil {
			return 
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// var url = "postgres://usher_admin:tiendada123@localhost/usher"
		var url = utils.GetFlagValue(cmd, "url", "")

		err := core.CreateDatabaseVersionTable(url)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
