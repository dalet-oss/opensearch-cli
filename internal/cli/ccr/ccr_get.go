package ccr

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var ccrGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "query remote settings for the cluster",
	Example: `opensearch-cli ccr get`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.NewFromCmd(cmd).
			GetRemoteSettings(flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
			log.Fatal().Msgf("failed to get remote cluster settings:%v", err)
		}
	},
}
