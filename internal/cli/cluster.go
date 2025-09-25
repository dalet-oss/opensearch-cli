package cli

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"fmt"
	"github.com/spf13/cobra"
)

func NewClusterCmd() *cobra.Command {
	// subcommands
	return clusterCmd
}

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"ctx"},
	Short:   "stub for tests | experiments",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		client := api.New(config)
		client.ClusterSettings()
		plugins := client.PluginsList()
		fmt.Println(api.HasPlugin(plugins, api.SecurityPlugin))
		fmt.Println(api.HasPlugin(plugins, api.CCRPlugin))
	},
}
