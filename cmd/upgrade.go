package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/giwiro/walkline/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use: "upgrade [flags] <target>",
	Example: `  walkline upgrade head
  walkline upgrade V002
`,
	Short:      "Upgrades database to the target version",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"revision"},
	Run: func(cmd *cobra.Command, args []string) {
		var firstVersion *core.VersionShort
		var targetVersion *core.VersionShort

		var verbose = utils.GetFlagBooleanValue(cmd, "verbose", false)
		var url = utils.GetFlagStringValue(cmd, "url", "")
		var path = utils.GetFlagStringValue(cmd, "path", "")

		if args[0] == "head" {
			targetVersion = nil
		} else {
			versionShort, err := core.ParseVersionShort(args[0])

			if err != nil {
				if verbose == true {
					log.Println("Bad version format: ", err)
				}
				os.Exit(1)
			}

			if versionShort.Prefix == "U" {
				if verbose == true {
					log.Println("Target version cannot be an undo migration: ", err)
				}
				os.Exit(1)
			}

			targetVersion = versionShort
		}

		firstNode, _, err := core.BuildMigrationTreeFromPath(path)

		if err != nil {
			if verbose == true {
				log.Println("Could not build migration tree:", err)
			}
			os.Exit(1)
		}

		currentVersion, flavor, err := core.GetCurrentDatabaseVersion(url, verbose)

		if err != nil {
			if verbose == true {
				fmt.Println("Could not get current DB version:", err)
			}
			os.Exit(1)
		}

		if currentVersion == nil {
			if verbose == true {
				fmt.Println("Found empty current DB version")
			}
			firstVersion = &core.VersionShort{
				Prefix:  firstNode.File.Version.Prefix,
				Version: firstNode.File.Version.Version,
			}
		} else {
			var currentNode = core.FindMigrationNode(firstNode, currentVersion)

			if currentNode != nil && currentNode.NextMigrationNode != nil {
				firstVersion = &core.VersionShort{
					Prefix:  currentNode.NextMigrationNode.File.Version.Prefix,
					Version: currentNode.NextMigrationNode.File.Version.Version,
				}
			} else {
				if verbose == true {
					log.Println("Could not generate first migration version")
				}
				os.Exit(1)
			}
		}

		migration, err := core.GenerateMigrationStringFromVersionShortRange(flavor, path, currentVersion, firstVersion, targetVersion)

		if err != nil {
			if verbose == true {
				log.Println("Could not generate migration:", err)
			}
			os.Exit(1)
		}

		err = core.ExecuteMigrationString(url, migration, verbose)

		if err != nil {
			if verbose == true {
				log.Println("Could not execute transaction: ", err)
			}
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
