package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"github.com/spf13/cobra"
	"log"
	"sort"
)

var statusReplicationCmd = &cobra.Command{
	Use:   "status",
	Short: "show replication status.",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.New(configutils.LoadConfig(flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag)))
		replicationIndex := ""
		if len(args) == 0 || args[0] == "" {
			log.Println("index name is required")
			registeredIndices := client.GetIndexList()
			indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string {
				return info.Index
			})
			sort.Strings(indexNames)
			replicationIndex = prompts.SelectivePrompt("Select index for query", indexNames)
			log.Println("selected index: ", replicationIndex)
		} else {
			replicationIndex = args[0]
		}
		client.StatusReplication(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), RawFlag))
	},
}

func init() {
	statusReplicationCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
