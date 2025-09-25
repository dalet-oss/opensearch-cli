/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package ccr

import "github.com/spf13/cobra"

func NewCCRCmd() *cobra.Command {
	ccrCmd.AddCommand(ccrCreateCmd)
	ccrCmd.AddCommand(ccrGetCmd)
	ccrCmd.AddCommand(ccrDeleteCmd)
	return ccrCmd
}

var ccrCmd = &cobra.Command{
	Use:   "ccr",
	Short: "cross-cluster replication settings management commands.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.HasAvailableSubCommands() {
			cmd.Help()
		}
	},
}
