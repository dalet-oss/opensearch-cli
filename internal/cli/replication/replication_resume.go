package replication

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"github.com/spf13/cobra"
	"sort"
)

var replicationResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "resume replication",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		replicationIndex := ""
		if len(args) == 0 || args[0] == "" {
			log.Info().Msg("index name is required")
			registeredIndices := client.GetIndexList()
			indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string {
				return info.Index
			})
			sort.Strings(indexNames)
			replicationIndex = prompts.SelectivePrompt("Select index for query", indexNames)
			log.Info().Msgf("selected index: %s", replicationIndex)
		} else {
			replicationIndex = args[0]
		}
		client.ResumeReplication(replicationIndex, flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))

	},
}
