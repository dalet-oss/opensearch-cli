package replication

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/prompts"
	"github.com/spf13/cobra"
	"sort"
)

var replicationResumeCmd = &cobra.Command{
	Use:     "resume",
	Short:   "resume replication",
	Example: `opensearch-cli replication resume [INDEX NAME | index pattern]`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		replicationIndex := ""
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
			log.Info().Msgf("selected index: %s", replicationIndex)
			if err := client.ResumeReplication(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
				log.Fatal().Err(err)
			}
		} else if !gu.ContainsWildcard(args[0]) {
			replicationIndex = args[0]
			if err := client.ResumeReplication(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
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
					log.Info().Msgf("resuming replication for index '%s'", index)
					if err := client.ResumeReplication(index, flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
						log.Fatal().Err(err)
					}
				}
			}
		}

	},
}
