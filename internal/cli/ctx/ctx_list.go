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

// ctxListCmd represents the list command
var ctxListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		config.ListContexts()
		fmt.Println()
		fmt.Println("Legend:\n\t✅ - entry found in the configuration file.\n\t❌ - entry not found in the configuration file. ")
	},
}
