package stats

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/prompts"
	"github.com/spf13/cobra"
	"sort"
)

const RawFlag = "raw"

var lagCmd = &cobra.Command{
	Use:     "lag",
	Short:   "show lag information for a specific index.",
	Example: `opensearch-cli stats lag [INDEX NAME | index pattern]`,
	Run: func(cmd *cobra.Command, args []string) {
		replicationIndex := ""
		client := api.NewFromCmd(cmd)
		if len(args) == 0 || args[0] == "" {
			log.Info().Msg("index name is required")
			registeredIndices, indexListErr := client.GetIndexList()
			if indexListErr != nil {
				log.Fatal().Err(indexListErr)
			}
			indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string {
				return info.Index
			})
			sort.Strings(indexNames)
			replicationIndex = prompts.SelectivePrompt("Select index for query", indexNames)
			fmt.Println("selected index: ", replicationIndex)
			if _, err := client.GetStatsLag(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), RawFlag)); err != nil {
				log.Fatal().Err(err)
			}
		} else if !gu.ContainsWildcard(args[0]) {
			replicationIndex = args[0]
			if _, err := client.GetStatsLag(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), RawFlag)); err != nil {
				log.Fatal().Err(err)
			}
		} else {
			registeredIndices, indexListErr := client.GetIndexList()
			if indexListErr != nil {
				log.Fatal().Err(indexListErr)
			}
			indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string { return info.Index })
			filtered := fp.Filter(indexNames, gu.GetMatchFunc(args[0]))
			sort.Strings(filtered)
			if len(filtered) == 0 {
				log.Warn().Msgf("no indices found for %s expression [total %d in the cluster]", args[0], len(indexNames))
				return
			} else {
				log.Info().Msgf(
					"found %d %s for %s expression",
					len(filtered), fp.Ternary("index", "indices", len(filtered) == 1), args[0])
				for _, index := range filtered {
					if _, err := client.GetStatsLag(index, flagutils.GetBoolFlag(cmd.Flags(), RawFlag)); err != nil {
						log.Fatal().Err(err)
					}
				}
			}
		}
	},
}
