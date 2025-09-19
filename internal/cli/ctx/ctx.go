/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package ctx

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCtxCmd() *cobra.Command {
	// subcommands
	ctxCmd.AddCommand(addCmd)
	ctxCmd.AddCommand(ctxCurrentCmd)
	ctxCmd.AddCommand(ctxListCmd)
	ctxCmd.AddCommand(ctxRemoveCmd)
	ctxCmd.AddCommand(ctxSwitchCmd)
	ctxCmd.AddCommand(ctxViewCmd)

	return ctxCmd
}

// ctxCmd represents the ctx command
var ctxCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ctx called")
	},
}
