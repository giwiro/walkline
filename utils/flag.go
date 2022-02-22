package utils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetFlagStringValue(cmd *cobra.Command, flag string, def string) string {
	var value string = def
	var viperValue = viper.GetString(flag)
	var f = cmd.Flag(flag)

	if len(f.Value.String()) > 0 {
		return f.Value.String()
	}

	if len(viperValue) > 0 {
		return viperValue
	}

	return value
}

func GetFlagBooleanValue(cmd *cobra.Command, flag string, def bool) bool {
	f, err := cmd.Flags().GetBool(flag)

	if err != nil {
		return def
	}

	return f || viper.GetBool(flag)
}
