/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package replication

import (
	"github.com/spf13/cobra"
)
import "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

func NewReplicationCmd() *cobra.Command {
	// subcommands
	replicationCmd.AddCommand(
		replicationCreateCmd,
		replicationPauseCmd,
		replicationStopCmd,
		replicationResumeCmd,
		replicationStatusCmd,
		replicationTaskStatusCmd,
	)
	return replicationCmd
}

// replicationCmd represents the ctx command
var replicationCmd = &cobra.Command{
	Use:     "replication",
	Aliases: []string{"repl"},
	Short:   "replication commands.",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			cmd.Help()
		}
	},
}
