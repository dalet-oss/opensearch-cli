package ctx

import (
	"github.com/spf13/cobra"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

func NewCtxCmd() *cobra.Command {
	// subcommands
	ctxCmd.AddCommand(
		addCmd,
		ctxCurrentCmd,
		ctxListCmd,
		ctxSwitchCmd,
		ctxViewCmd,
	)

	return ctxCmd
}

// ctxCmd represents the ctx command
var ctxCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"},
	Short:   "manage contexts, clusters and users.",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.HasAvailableSubCommands() {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
		}
	},
}
