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

const (
	DetailedFlag = "detailed"
)

var replicationTaskStatusCmd = &cobra.Command{
	Use:     "task-status",
	Short:   "show replication task status",
	Example: `opensearch-cli replication task-status [INDEX NAME | index pattern] [--detailed]`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		replicationIndex := ""
		detailed := flagutils.GetBoolFlag(cmd.Flags(), DetailedFlag)
		raw := flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)
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
			if err := client.TaskStatusReplication(replicationIndex, detailed, raw); err != nil {
				log.Fatal().Err(err)
			}
		} else if !gu.ContainsWildcard(args[0]) {
			replicationIndex = args[0]
			if err := client.TaskStatusReplication(replicationIndex, detailed, raw); err != nil {
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
				log.Info().Msgf("found %d indices for '%s' expression", len(filtered), args[0])
				for _, index := range filtered {
					log.Info().Msgf("querying replication status for the index '%s'", index)
					if err := client.TaskStatusReplication(index, detailed, raw); err != nil {
						log.Fatal().Err(err)
					}
				}
			}
		}
	},
}

func init() {
	replicationTaskStatusCmd.PersistentFlags().Bool(DetailedFlag, false, "show detailed info about tasks.")
}
