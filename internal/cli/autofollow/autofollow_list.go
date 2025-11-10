package autofollow

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var autofollowListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "shows list of configured autofollow rules.",
	Example: "opensearch-cli autofollow list",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := api.NewFromCmd(cmd).
			ListOfAFRules(flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
			log.Fatal().Msgf("failed to get the list autofollow rules:%v", err)
		}
	},
}
