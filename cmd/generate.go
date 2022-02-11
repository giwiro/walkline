package cmd

import (
	"fmt"
	"github.com/giwiro/walkline/core"
	"github.com/spf13/cobra"
	"log"
)


// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:        "generate [flags] <revision_range>",
	Example:    `  walkline generate V001:V002
  walkline generate U001:U001
  walkline generate --flavor=postgresql U001:U001
`,
	Short:      "Generates sql revision based on the version ranged provided",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"revision_range"},
	Run: func(cmd *cobra.Command, args []string) {
		leftVersion, rightVersion, err := core.ParseVersionShortRange(args[0])

		if err != nil {
			log.Fatal("Could not parse version range")
		}

		var flavor string
		var flavorFlag = cmd.Flag("flavor")

		fmt.Println("flavorFlag", flavorFlag.Value)

		if len(flavorFlag.Value.String()) > 0 {
			flavor = flavorFlag.Value.String()
		}else {
			flavor = "postgresql"
		}

		transaction, err := core.GenerateMigrationStringFromVersionShortRange(flavor, leftVersion, rightVersion)
		if err != nil {
			fmt.Println(err)
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
