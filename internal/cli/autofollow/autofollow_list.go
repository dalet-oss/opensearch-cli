package autofollow

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var autofollowListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "shows list of configured autofollow rules.",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).
			ListOfAFRules(
				flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag),
			)
	},
}
