package replication

import (
	"fmt"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/spf13/cobra"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

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
	Long: fmt.Sprintf(
		`
Some of the commands support wildcard operations:
%s`, gu.WildHelp),
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
		}
	},
}
