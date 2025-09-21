/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package ctx

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"fmt"
	"github.com/spf13/cobra"
)

// ctxCurrentCmd represents the current command
var ctxCurrentCmd = &cobra.Command{
	Use:   "show",
	Short: "show active context information.",
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		fmt.Println(config.ShowContextInfo(config.Current))
		fmt.Println()
		fmt.Println("Legend:\n\t✅ - entry found in the configuration file.\n\t❌ - entry not found in the configuration file. ")
	},
}
