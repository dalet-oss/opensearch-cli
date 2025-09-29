package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"sort"
)

const RawFlag = "raw"

var lagCmd = &cobra.Command{
	Use:   "lag",
	Short: "show lag information for a specific index.",
	Run: func(cmd *cobra.Command, args []string) {
		replicationIndex := ""
		client := api.NewFromCmd(cmd)
		if len(args) == 0 || args[0] == "" {
			log.Println("index name is required")
			registeredIndices := client.GetIndexList()
			indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string {
				return info.Index
			})
			sort.Strings(indexNames)
			replicationIndex = prompts.SelectivePrompt("Select index for query", indexNames)
			fmt.Println("selected index: ", replicationIndex)
		} else {
			replicationIndex = args[0]
		}
		client.GetStatsLag(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), RawFlag))
	},
}

func init() {
	lagCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
