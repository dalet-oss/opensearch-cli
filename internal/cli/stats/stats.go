/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package stats

import (
	"github.com/spf13/cobra"
)

func NewStatsCmd() *cobra.Command {
	// subcommands
	statsCmd.AddCommand(
		lagCmd,
		autofollowRstatsCmd,
		leaderRStatsCmd,
		followerRStatsCmd,
	)
	return statsCmd
}

// ctxCmd represents the ctx command
var statsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"st"},
	Short:   "show stats information.",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			cmd.Help()
		}
	},
}
