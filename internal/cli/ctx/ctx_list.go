/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
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
	Short: "list all contexts.",
	Long:  `Show information about all contexts.`,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		config.ListContexts()
		fmt.Println()
		fmt.Println("Legend:\n\t✅ - entry found in the configuration file.\n\t❌ - entry not found in the configuration file. ")
	},
}
