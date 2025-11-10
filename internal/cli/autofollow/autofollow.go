package autofollow

import "github.com/spf13/cobra"
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

const (
	IndexPatternFlag        = "pattern"
	LeaderAliasFlag         = "leader"
	LeaderClusterRoleFlag   = "leader-cluster-role"
	FollowerClusterRoleFlag = "follower-cluster-role"
)

func NewAutofollowCmd() *cobra.Command {
	autofollowCmd.AddCommand(
		autofollowCreateCmd,
		autofollowDeleteCmd,
		autofollowListCmd,
	)
	return autofollowCmd
}

var autofollowCmd = &cobra.Command{
	Use:     "autofollow",
	Aliases: []string{"af"},
	Short:   "Manage autofollow settings for the OpenSearch cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
		}
	},
}
