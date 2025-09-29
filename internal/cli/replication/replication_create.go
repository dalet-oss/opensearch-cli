package replication

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/replication"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
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
	Use:   "create",
	Short: "start replication thing",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		options := prepareReplicationCall(cmd.Flags(), client)
		client.CreateReplication(options, flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))
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
	if api.HasPlugin(client.PluginsList(), api.SecurityPlugin) {
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
