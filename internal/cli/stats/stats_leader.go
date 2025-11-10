package stats

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var leaderRStatsCmd = &cobra.Command{
	Use:     "leader",
	Short:   "show leader replication stats.",
	Example: `opensearch-cli stats leader`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := api.NewFromCmd(cmd).GetReplicationLeaderStats(flagutils.GetBoolFlag(cmd.Flags(), RawFlag)); err != nil {
			log.Fatal().Msgf("failed to get leader stats:%v", err)
		}
	},
}
