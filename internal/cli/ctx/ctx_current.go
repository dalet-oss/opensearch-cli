/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package ctx

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	configutils "github.com/dalet-oss/opensearch-cli/pkg/utils/config"
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
