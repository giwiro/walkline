package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/giwiro/walkline/utils"
	"github.com/spf13/cobra"
	"log"
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

		var url = utils.GetFlagValue(cmd, "url", "")
		var path = utils.GetFlagValue(cmd, "path", "")

		if args[0] == "head" {
			targetVersion = nil
		} else {
			versionShort, err := core.ParseVersionShort(args[0])

			if err != nil {
				log.Fatal("Bad version format: ", err)
			}

			if versionShort.Prefix == "U" {
				log.Fatal("Target version cannot be an undo migration")
			}

			targetVersion = versionShort
		}

		firstNode, _, err := core.BuildMigrationTreeFromPath(path)

		if err != nil {
			log.Fatal(err)
		}

		currentVersion, flavor, err := core.GetCurrentDatabaseVersion(url)

		if err != nil {
			fmt.Println("Could not get current DB version: ", err)
			firstVersion = &core.VersionShort{
				Prefix: firstNode.File.Version.Prefix,
				Version: firstNode.File.Version.Version,
			}
		} else {
			var currenNode = core.FindMigrationNode(firstNode, currentVersion)

			if currenNode == nil || currenNode.NextMigrationNode == nil {
				firstVersion =  &core.VersionShort{
					Prefix: currenNode.NextMigrationNode.File.Version.Prefix,
					Version: currenNode.NextMigrationNode.File.Version.Version,
				}
			} else {
				log.Fatal("Could not generate first migration version")
			}
		}

		migration, err := core.GenerateMigrationStringFromVersionShortRange(flavor, path, currentVersion, firstVersion, targetVersion)

		if err != nil {
			log.Fatal("Could not generate migration: ", err)
		}

		err = core.ExecuteMigrationString(url, migration)

		if err != nil {
			log.Fatal("Could not execute transaction: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")
	upgradeCmd.PersistentFlags().String("path", "", "Path of the migration files")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
