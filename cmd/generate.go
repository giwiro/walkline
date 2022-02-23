package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/giwiro/walkline/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use: "generate [flags] <revision_range>",
	Example: `  walkline generate V001:V002
  walkline generate U001:U001 (This will generate just the U001)
  walkline generate U001 (This will generate just the U001 as well)
  walkline generate V002 (This will generate all revisions until V002)
  walkline generate --flavor=postgresql U001:U001
`,
	Short:      "Generates sql revision based on the version ranged provided",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"revision_range"},
	Run: func(cmd *cobra.Command, args []string) {
		var leftVersion *core.VersionShort
		var rightVersion *core.VersionShort

		var path = utils.GetFlagStringValue(cmd, "path", "")
		var verbose = utils.GetFlagBooleanValue(cmd, "verbose", false)
		var schema = utils.GetFlagStringValue(cmd, "schema", "")
		var url = utils.GetFlagStringValue(cmd, "url", "")

		_, flavor, err := core.GetCurrentDatabaseVersion(url, verbose, schema)

		singleVersion, err := core.ParseVersionShort(args[0])

		firstNode, _, buildTreeErr := core.BuildMigrationTreeFromPath(path)

		if buildTreeErr != nil {
			if verbose == true {
				log.Println("Could not build migration tree:", buildTreeErr)
			}
			os.Exit(1)
		}

		if firstNode == nil {
			if verbose == true {
				log.Println("Could not found first node")
			}
			os.Exit(1)
		}

		if err == nil {
			if singleVersion.Prefix == "U" {
				leftVersion = singleVersion
				rightVersion = singleVersion
			} else {
				leftVersion = core.GetVersionShortFromFull(firstNode.File.Version)
				rightVersion = singleVersion
			}
		} else {
			leftVersion, rightVersion, err = core.ParseVersionShortRange(args[0])

			if err != nil {
				if verbose == true {
					log.Println("Could not parse version range:", err)
				}
				os.Exit(1)
			}
		}

		transaction, err := core.GenerateMigrationStringFromVersionShortRange(flavor, path, schema, leftVersion, leftVersion, rightVersion)

		if err != nil {
			if verbose == true {
				log.Println("Could not generate migration string:", err)
			}
		}

		fmt.Println(transaction)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
