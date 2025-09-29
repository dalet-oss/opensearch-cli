package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var followerRStatsCmd = &cobra.Command{
	Use:   "follower",
	Short: "show follower replication stats.",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).
			GetReplicationFollowerStats(
				flagutils.GetBoolFlag(cmd.Flags(), RawFlag),
			)
	},
}

func init() {
	followerRStatsCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
