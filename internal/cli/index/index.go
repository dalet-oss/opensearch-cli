package index

import "github.com/spf13/cobra"

func NewIndexCmd() *cobra.Command {
	// subcommands
	indexCmd.AddCommand(indexListCmd)
	indexCmd.AddCommand(indexDeleteCmd)
	return indexCmd
}

// ctxCmd represents the ctx command
var indexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"index"},
	Short:   "index commands",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			cmd.Help()
		}
	},
}
