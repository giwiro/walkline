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
        var schema = utils.GetFlagStringValue(cmd, "schema", "")

        versionShort, _, err := core.GetCurrentDatabaseVersion(url, verbose, schema)

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
}
