package stats

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var followerRStatsCmd = &cobra.Command{
	Use:     "follower",
	Short:   "show follower replication stats.",
	Example: `opensearch-cli stats follower`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := api.NewFromCmd(cmd).GetReplicationFollowerStats(flagutils.GetBoolFlag(cmd.Flags(), RawFlag)); err != nil {
			log.Fatal().Msgf("failed to get follower stats:%v", err)
		}
	},
}
