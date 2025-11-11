package ctx

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	configutils "github.com/dalet-oss/opensearch-cli/pkg/utils/config"
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
