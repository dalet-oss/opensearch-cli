/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package stats

import (
	"github.com/spf13/cobra"
)

func NewReplicationCmd() *cobra.Command {
	// subcommands
	replicationCmd.AddCommand()
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
