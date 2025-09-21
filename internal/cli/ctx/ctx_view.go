/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package ctx

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"fmt"
	"github.com/spf13/cobra"
)

// ctxViewCmd represents the view command
var ctxViewCmd = &cobra.Command{
	Use:   "view",
	Short: "show active context information.",
	Long:  `Show entire information about the active context(except the credentials)`,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		var info string
		if len(args) > 0 {
			info = config.ShowContextInfoExtended(args[0])
		} else {
			info = config.ShowContextInfoExtended(config.Current)
		}
		fmt.Println(info)
	},
}
