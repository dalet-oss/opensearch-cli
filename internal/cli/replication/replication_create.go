package replication

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	IndexNameFlag           = "index"
	LeaderAliasFlag         = "leader"
	LeaderIndexFlag         = "leader-index"
	LeaderClusterRoleFlag   = "cluster-role"
	FollowerClusterRoleFlag = "follower-cluster-role"
)

var replicationCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create index replication task",
	Example: `opensearch-cli replication create --index <INDEX NAME> --leader leader --leader-index [--cluster-role leader-role] [--follower-cluster-role follower-role]`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		options := prepareReplicationCall(cmd.Flags(), client)
		if err := client.CreateReplication(options, flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
			log.Fatal().Msgf("failed to create replication task:%v", err)
		}
	},
}

// prepareReplicationCall gather options from cli required args to start replication
// TODO: interactive mode if no args supplied
func prepareReplicationCall(flags *pflag.FlagSet, client *api.OpensearchWrapper) replication.StartReplicationReq {
	opts := replication.StartReplicationReq{
		Index: flagutils.GetNotEmptyStringFlag(flags, IndexNameFlag),
		Body: replication.StartReplicationBody{
			LeaderAlias: flagutils.GetNotEmptyStringFlag(flags, LeaderAliasFlag),
			LeaderIndex: flagutils.GetNotEmptyStringFlag(flags, LeaderIndexFlag),
		},
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
	replicationCreateCmd.PersistentFlags().String(IndexNameFlag, "", "name of the index to create replication")
	replicationCreateCmd.PersistentFlags().String(LeaderAliasFlag, "", "leader alias")
	replicationCreateCmd.PersistentFlags().String(LeaderIndexFlag, "", "leader index")
	replicationCreateCmd.PersistentFlags().String(LeaderClusterRoleFlag, "", "[mandatory if security plugin enabled]leader cluster role")
	replicationCreateCmd.PersistentFlags().String(FollowerClusterRoleFlag, "", "[mandatory if security plugin enabled]follower cluster role")
}
