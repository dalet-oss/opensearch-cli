/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package ctx

import (
	"github.com/spf13/cobra"
)

func NewCtxCmd() *cobra.Command {
	// subcommands
	ctxCmd.AddCommand(addCmd)
	ctxCmd.AddCommand(ctxCurrentCmd)
	ctxCmd.AddCommand(ctxListCmd)
	ctxCmd.AddCommand(ctxSwitchCmd)
	ctxCmd.AddCommand(ctxViewCmd)

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
