package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/giwiro/walkline/utils"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var downgradeTimesRegex = regexp.MustCompile("^\\d+?$")

// downgradeCmd represents the downgrade command
var downgradeCmd = &cobra.Command{
	Use:     "downgrade [flags] <times>",
	Example: "  walkline downgrade 1",
	Short:   "Downgrades database n times",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var verbose = utils.GetFlagBooleanValue(cmd, "verbose", false)
		var url = utils.GetFlagStringValue(cmd, "url", "")
		var path = utils.GetFlagStringValue(cmd, "path", "")
		var schema = utils.GetFlagStringValue(cmd, "schema", "")

		if !downgradeTimesRegex.MatchString(args[0]) {
			fmt.Println("Downgrade times must be a number")
			os.Exit(1)
		}

		times, err := strconv.Atoi(args[0])

		if err != nil {
			fmt.Println("Could not parse downgrade times")
			os.Exit(1)
		}

		currentVersion, flavor, err := core.GetCurrentDatabaseVersion(url, verbose, schema)

		if err != nil {
			if verbose == true {
				log.Println("Could not get database version:", err)
			}
			os.Exit(1)
		}

		if currentVersion == nil {
			if verbose == true {
				log.Println("Could not get database version: version empty")
			}
			os.Exit(1)
		}

		migration, err := core.GenerateConsecutiveDowngradesMigrationString(flavor, path, schema, currentVersion, times)

		if err != nil {
			if verbose == true {
				log.Println("Could not build migration tree:", err)
			}
			os.Exit(1)
		}

		err = core.ExecuteMigrationString(url, migration, verbose)

		if err != nil {
			if verbose == true {
				log.Println("Could not execute transaction:", err)
			}
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
