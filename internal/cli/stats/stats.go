/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package stats

import (
	"github.com/spf13/cobra"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

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
	Short:   "Collection of commands showing stats information.",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
		}
	},
}
