package utils

import (
	"github.com/spf13/cobra"
)

func GetFlagValue(cmd *cobra.Command, flag string, def string) string {
	var value string = def
	var f = cmd.Flag(flag)

	if len(f.Value.String()) > 0 {
		return f.Value.String()
	}

	return value
}
