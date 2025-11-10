package index

import "github.com/spf13/cobra"
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

func NewIndexCmd() *cobra.Command {
	// subcommands
	indexCmd.AddCommand(
		indexListCmd,
		indexDeleteCmd,
		indexCreateCmd,
	)
	return indexCmd
}

// ctxCmd represents the ctx command
var indexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"index"},
	Short:   "index commands",
	Long:    `Set of commands for index management`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
		}
	},
}

func init() {
	indexDeleteCmd.Flags().Bool(ConfirmFlag, false, "delete index without confirmation")
	indexListCmd.Flags().Bool(FlagAll, false, "show all indices, including hidden ones(starting with '.').")
}
