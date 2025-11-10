package stats

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var autofollowRstatsCmd = &cobra.Command{
	Use:     "autofollow",
	Short:   "show autofollow information.",
	Example: `opensearch-cli stats autofollow`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.NewFromCmd(cmd).GetReplicationAutofollowStats(flagutils.GetBoolFlag(cmd.Flags(), RawFlag)); err != nil {
			log.Fatal().Msgf("failed to get autofollow stats:%v", err)
		}
	},
}
