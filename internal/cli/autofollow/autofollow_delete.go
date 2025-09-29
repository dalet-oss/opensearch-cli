package autofollow

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/replication"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var autofollowDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete autofollow rule from the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		client.
			DeleteAutofollow(autofollowDeleteOpts(cmd.Flags(), client), flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))
	},
}

func autofollowDeleteOpts(flags *pflag.FlagSet, client *api.OpensearchWrapper) replication.DeleteAutofollowReq {
	opts := replication.DeleteAutofollowReq{}
	opts.Body = replication.DeleteAutofollowBody{
		Name:        flagutils.GetNotEmptyStringFlag(flags, RuleNameFlag),
		LeaderAlias: flagutils.GetNotEmptyStringFlag(flags, LeaderAliasFlag),
	}
	return opts
}

func init() {
	autofollowDeleteCmd.PersistentFlags().String(RuleNameFlag, "", "name of the autofollow rule")
	autofollowDeleteCmd.PersistentFlags().String(LeaderAliasFlag, "", "leader alias")
}
