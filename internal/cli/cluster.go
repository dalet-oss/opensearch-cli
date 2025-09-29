package cli

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
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
		client := api.NewFromCmd(cmd)
		client.ClusterSettings()
		plugins := client.PluginsList()
		fmt.Println(api.HasPlugin(plugins, api.SecurityPlugin))
		fmt.Println(api.HasPlugin(plugins, api.CCRPlugin))
	},
}
