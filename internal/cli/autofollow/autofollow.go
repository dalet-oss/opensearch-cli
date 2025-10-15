package autofollow

import "github.com/spf13/cobra"

const (
	RuleNameFlag            = "name"
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
	Short:   "management autofollow settings for replication",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			cmd.Help()
		}
	},
}
