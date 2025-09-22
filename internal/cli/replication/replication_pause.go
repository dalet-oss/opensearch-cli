package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"github.com/spf13/cobra"
)

var replicationPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause replication",
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		client := api.New(configutils.LoadConfig(appConfigFile))
	},
}

func init() {
	replicationPauseCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
