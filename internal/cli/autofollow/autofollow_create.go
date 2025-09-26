package autofollow

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/replication"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var autofollowCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"update"},
	Short:   "Create or update autofollow rule in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.New(
			configutils.LoadConfig(
				flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag)))
		client.
			CreateAutofollowRule(prepareAutofollowOpts(cmd.Flags(), client), flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))
	},
}

func prepareAutofollowOpts(flags *pflag.FlagSet, client *api.OpensearchWrapper) replication.CreateAutofollowReq {
	opts := replication.CreateAutofollowReq{}
	opts.Body = replication.CreateAutofollowBody{
		Name:         flagutils.GetNotEmptyStringFlag(flags, RuleNameFlag),
		LeaderAlias:  flagutils.GetNotEmptyStringFlag(flags, LeaderAliasFlag),
		IndexPattern: flagutils.GetNotEmptyStringFlag(flags, IndexPatternFlag),
	}
	if api.HasPlugin(client.PluginsList(), api.SecurityPlugin) {
		opts.Body.UseRoles = replication.ReplicationRoles{
			LeaderClusterRole:   flagutils.GetNotEmptyStringFlag(flags, LeaderClusterRoleFlag),
			FollowerClusterRole: flagutils.GetNotEmptyStringFlag(flags, FollowerClusterRoleFlag),
		}
	}
	return opts
}

func init() {
	autofollowCreateCmd.PersistentFlags().String(RuleNameFlag, "", "name of the autofollow rule")
	autofollowCreateCmd.PersistentFlags().String(LeaderAliasFlag, "", "leader alias[could be created with opensearch ccr create]")
	autofollowCreateCmd.PersistentFlags().String(IndexPatternFlag, "", "index pattern of the autofollow rule")
	autofollowCreateCmd.PersistentFlags().String(LeaderClusterRoleFlag, "", "[mandatory if security plugin enabled]leader cluster role")
	autofollowCreateCmd.PersistentFlags().String(FollowerClusterRoleFlag, "", "[mandatory if security plugin enabled]follower cluster role")
}
