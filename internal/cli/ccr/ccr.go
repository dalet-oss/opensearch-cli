package ccr

import (
	"github.com/dalet-oss/opensearch-cli/pkg/utils/logging"
	"github.com/spf13/cobra"
)

var log = logging.Logger()

func NewCCRCmd() *cobra.Command {
	ccrCmd.AddCommand(
		ccrCreateCmd,
		ccrGetCmd,
		ccrDeleteCmd,
	)
	return ccrCmd
}

var ccrCmd = &cobra.Command{
	Use:   "ccr",
	Short: "cross-cluster replication settings management commands.",
	Long:  `Set of commands for management cross-cluster replication.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
		}
	},
}
