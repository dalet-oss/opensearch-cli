package ccr

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var ccrGetCmd = &cobra.Command{
	Use:   "get",
	Short: "query remote settings for the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).
			GetRemoteSettings(flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))
	},
}
