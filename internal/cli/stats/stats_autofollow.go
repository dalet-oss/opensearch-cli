package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var autofollowRstatsCmd = &cobra.Command{
	Use:   "autofollow",
	Short: "show autofollow information.",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).
			GetReplicationAutofollowStats(
				flagutils.GetBoolFlag(cmd.Flags(), RawFlag),
			)
	},
}

func init() {
	autofollowRstatsCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
