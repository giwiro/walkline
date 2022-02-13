package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"log"

	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use: 		"upgrade [flags] <target>",
	Example: 	`  walkline upgrade head
  walkline upgrade V002
`,
	Short: 		"Upgrades database to the target version",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"revision"},
	Run: func(cmd *cobra.Command, args []string) {
		var targetVersion *core.VersionShort
		var url string
		var urlFlag = cmd.Flag("url")

		fmt.Println("urlFlag", urlFlag.Value)

		if len(urlFlag.Value.String()) > 0 {
			url = urlFlag.Value.String()
		}

		if args[0] == "head" {
			targetVersion = nil
		} else {
			versionShort, err := core.ParseVersionShort(args[0])

			if err != nil {
				log.Fatal(err)
			}

			targetVersion = versionShort
		}

		firstNode, _, err := core.BuildMigrationTreeFromPath(migrationPath)

		if err != nil {
			log.Fatal(err)
		}

		version, flavor, err := core.GetCurrentDatabaseVersion(url)

		if err != nil {
			fmt.Println("Could not get current DB version")
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
