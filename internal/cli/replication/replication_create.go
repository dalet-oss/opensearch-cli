package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/replication"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
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
		client := api.New(
			configutils.LoadConfig(
				flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag)),
		)
		options := prepareReplicationCall(cmd.Flags(), client)
		client.CreateReplication(options, flagutils.GetBoolFlag(cmd.Flags(), RawFlag))
	},
}

// prepareReplicationCall gather options from cli required args to start replication
// TODO: interactive mode if no args supplied
func prepareReplicationCall(flags *pflag.FlagSet, client *api.OpensearchWrapper) replication.StartReplicationReq {
	opts := replication.StartReplicationReq{}
	if indexName := flagutils.GetStringFlag(flags, IndexNameFlag); indexName != "" {
		opts.Index = indexName
	} else {
		log.Fatalf("--%s is required", IndexNameFlag)
	}
	if leaderName := flagutils.GetStringFlag(flags, LeaderAliasFlag); leaderName != "" {
		opts.Body = replication.StartReplicationBody{LeaderAlias: leaderName}
	} else {
		log.Fatalf("--%s is required", LeaderAliasFlag)
	}
	if leaderIndex := flagutils.GetStringFlag(flags, LeaderIndexFlag); leaderIndex != "" {
		opts.Body.LeaderIndex = leaderIndex
	} else {
		log.Fatalf("--%s is required", LeaderIndexFlag)
	}
	if api.HasPlugin(client.PluginsList(), api.SecurityPlugin) {
		if leaderRole := flagutils.GetStringFlag(flags, LeaderClusterRoleFlag); leaderRole != "" {
			opts.Body.UseRoles = replication.ReplicationRoles{LeaderClusterRole: leaderRole}
		} else {
			log.Fatalf("--%s is required because security plugin enabled", LeaderClusterRoleFlag)
		}
		if followerRole := flagutils.GetStringFlag(flags, FollowerClusterRoleFlag); followerRole != "" {
			opts.Body.UseRoles.FollowerClusterRole = followerRole
		} else {
			log.Fatalf("--%s is required because security plugin enabled", FollowerClusterRoleFlag)
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
	replicationCreateCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
