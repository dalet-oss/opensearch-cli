package ccr

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var ccrDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "delete remote configuration from the OpenSearch cluster",
	Example: `opensearch-cli ccr delete <NAME>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			log.Fatal().Msg("config name is required")
		}
		if err := api.NewFromCmd(cmd).
			DeleteRemote(args[0], flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
			log.Fatal().Msgf("failed to delete remote cluster:%v", err)
		}
	},
}
