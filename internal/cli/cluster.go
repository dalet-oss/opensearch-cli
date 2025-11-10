package cli

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/rs/zerolog/log"
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
	Hidden:  true,
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		client.ClusterSettings()
		plugins, queryPluginErr := client.PluginsList()
		if queryPluginErr != nil {
			log.Fatal().Msgf("fail to get plugin list:%v", queryPluginErr)
		}
		fmt.Println(api.HasPlugin(plugins, api.SecurityPlugin))
		fmt.Println(api.HasPlugin(plugins, api.CCRPlugin))
	},
}
