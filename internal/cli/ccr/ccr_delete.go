package ccr

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

var ccrDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del"},
	Short:   "delete remote configuration from the OpenSearch cluster",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).
			DeleteRemote(
				flagutils.GetNotEmptyStringFlag(cmd.Flags(), SettingsRemoteNameFlag),
				flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag),
			)
	},
}

func init() {
	ccrDeleteCmd.PersistentFlags().String(SettingsRemoteNameFlag, "", "name of the remote to delete")
}
