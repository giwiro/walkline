package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/giwiro/walkline/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var migrationPath = utils.GetFlagStringValue(cmd, "path", "")

		var verbose = utils.GetFlagBooleanValue(cmd, "verbose", false)
		var url = utils.GetFlagStringValue(cmd, "url", "")

		versionShort, _, err := core.GetCurrentDatabaseVersion(url, verbose)

		if (versionShort == nil || err != nil) && verbose == true {
			fmt.Println("Could not get current DB version:", err)
		}

		firstNode, _, err := core.BuildMigrationTreeFromPath(migrationPath)

		if err != nil {
			if verbose == true {
				log.Println("Could not build migration tree:", err)
			}
			os.Exit(1)
		}
		// fmt.Println(firstNode.NextMigrationNode.File.Content)
		core.PrintMigrationTree(firstNode, versionShort)
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// historyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// historyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
