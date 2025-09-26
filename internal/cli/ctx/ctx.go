/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package ctx

import (
	"github.com/spf13/cobra"
)

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
			cmd.Help()
		}
	},
}
