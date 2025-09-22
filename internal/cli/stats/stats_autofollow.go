package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"github.com/spf13/cobra"
)

var autofollowCmd = &cobra.Command{
	Use:   "autofollow",
	Short: "show autofollow information.",
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		_ = api.New(configutils.LoadConfig(appConfigFile))
		//client.GetStatsLag()
	},
}
