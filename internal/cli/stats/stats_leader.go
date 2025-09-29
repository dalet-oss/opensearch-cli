package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var leaderRStatsCmd = &cobra.Command{
	Use:   "leader",
	Short: "show leader replication stats.",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).GetReplicationLeaderStats(
			flagutils.GetBoolFlag(cmd.Flags(), RawFlag),
		)
	},
}

func init() {
	leaderRStatsCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
