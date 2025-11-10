package autofollow

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var autofollowCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"update"},
	Short:   "Create or update autofollow rule in the cluster",
	Long:    "Create | update autofollow rule in the cluster",
	Example: `autofollow create <RULE NAME> -l leader -p index-pattern [-r leader-role] [-f follower-role]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
			return
		}
		client := api.NewFromCmd(cmd)
		if err := client.
			CreateAutofollowRule(prepareAutofollowOpts(cmd.Flags(), args[0], client), flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
			log.Fatal().Msgf("failed to create autofollow rule:%v", err)
		}
	},
}

func prepareAutofollowOpts(flags *pflag.FlagSet, name string, client *api.OpensearchWrapper) replication.CreateAutofollowReq {
	opts := replication.CreateAutofollowReq{}
	opts.Body = replication.CreateAutofollowBody{
		Name:         name,
		LeaderAlias:  flagutils.GetNotEmptyStringFlag(flags, LeaderAliasFlag),
		IndexPattern: flagutils.GetNotEmptyStringFlag(flags, IndexPatternFlag),
	}
	plugins, queryPluginErr := client.PluginsList()
	if queryPluginErr != nil {
		log.Fatal().Msgf("fail to get plugin list:%v", queryPluginErr)
	}
	if api.HasPlugin(plugins, api.SecurityPlugin) {
		opts.Body.UseRoles = replication.ReplicationRoles{
			LeaderClusterRole:   flagutils.GetNotEmptyStringFlag(flags, LeaderClusterRoleFlag),
			FollowerClusterRole: flagutils.GetNotEmptyStringFlag(flags, FollowerClusterRoleFlag),
		}
	}
	return opts
}

func init() {
	//autofollowCreateCmd.PersistentFlags().String(RuleNameFlag, "", "name of the autofollow rule")
	autofollowCreateCmd.PersistentFlags().String(LeaderAliasFlag, "", "leader alias[could be created with opensearch ccr create]")
	autofollowCreateCmd.PersistentFlags().String(IndexPatternFlag, "", "index pattern of the autofollow rule")
	autofollowCreateCmd.PersistentFlags().String(LeaderClusterRoleFlag, "", "[mandatory if security plugin enabled]leader cluster role")
	autofollowCreateCmd.PersistentFlags().String(FollowerClusterRoleFlag, "", "[mandatory if security plugin enabled]follower cluster role")
}
