package replication

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
)

const (
	DetailedFlag = "detailed"
	TableFlag    = "table"
)

var replicationTaskStatusCmd = &cobra.Command{
	Use:   "task-status",
	Short: "show replication task status",
	Run: func(cmd *cobra.Command, args []string) {
		api.NewFromCmd(cmd).TaskStatusReplication(
			flagutils.GetStringFlag(cmd.Flags(), IndexNameFlag),
			flagutils.GetBoolFlag(cmd.Flags(), DetailedFlag),
			flagutils.GetBoolFlag(cmd.Flags(), TableFlag),
			flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))
	},
}

func init() {
	replicationTaskStatusCmd.PersistentFlags().Bool(DetailedFlag, false, "show detailed info about tasks.")
	replicationTaskStatusCmd.PersistentFlags().Bool(TableFlag, false, "show info as table")
	replicationTaskStatusCmd.PersistentFlags().String(IndexNameFlag, "", "show tasks for the index.")
}
