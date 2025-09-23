package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"github.com/spf13/cobra"
)

var replicationTaskStatusCmd = &cobra.Command{
	Use:   "task-status",
	Short: "show leader replication stats.",
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		client := api.New(configutils.LoadConfig(appConfigFile))
		showRawResp, _ := cmd.Flags().GetBool(RawFlag)
		if showRawResp {
		}
		client.TaskStatusReplication()
	},
}

func init() {
	replicationTaskStatusCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
