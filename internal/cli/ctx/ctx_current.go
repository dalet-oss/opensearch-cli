/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package ctx

import (
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"fmt"
	"github.com/spf13/cobra"
)

// ctxCurrentCmd represents the current command
var ctxCurrentCmd = &cobra.Command{
	Use:   "show",
	Short: "show active context information.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := configutils.LoadConfig("")
		fmt.Println(config.ShowContextInfo(config.Current))
	},
}
