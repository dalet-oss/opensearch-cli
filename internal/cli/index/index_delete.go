package index

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/prompts"
	"github.com/spf13/cobra"
	"slices"
	"sort"
	"strings"
)

const ConfirmFlag = "approve"

var indexDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "⚠️deletes index.",
	Long:  `Delete index in the OpenSearch cluster`,
	Example: fmt.Sprintf(`
opensearch-cli index delete <- will show interactive prompt with the list of indices
opensearch-cli index delete index1 <- will delete index1 from the OpenSearch cluster
opensearch-cli index delete index* <- will delete indices compliant with the pattern from the OpenSearch cluster
%s
`, gu.WildHelp),
	Run: func(cmd *cobra.Command, args []string) {
		// method
		client := api.NewFromCmd(cmd)
		indexToDelete := ""
		registeredIndices, indexListErr := client.GetIndexList()
		if indexListErr != nil {
			log.Fatal().Err(indexListErr)
		}
		indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string {
			return info.Index
		})
		sort.Strings(indexNames)
		if len(args) == 0 {
			indexToDelete = prompts.SelectivePrompt("Select index for removal", indexNames)
		} else if !gu.ContainsWildcard(args[0]) {
			// check if there's index -> delete
			indexToDelete = args[0]
			if idx := slices.IndexFunc(registeredIndices, func(info api.IndexInfo) bool {
				return info.Index == indexToDelete
			}); idx == -1 {
				log.Fatal().Msgf("❌index '%s' not found", indexToDelete)
			}
			if flagutils.GetBoolFlag(cmd.Flags(), ConfirmFlag) ||
				prompts.IsOk(
					prompts.QuestionPrompt(
						fmt.Sprintf(
							"[context:%s]Are you sure you want to delete index '%s'?", client.Config.Current,
							indexToDelete))) {
				indexDeleteErr := client.DeleteIndex(indexToDelete)
				if indexDeleteErr != nil {
					log.Fatal().Err(indexDeleteErr)
				}
			}
		} else {
			filtered := fp.Filter(indexNames, gu.GetMatchFunc(args[0]))
			sort.Strings(filtered)
			if len(filtered) == 0 {
				log.Warn().Msgf("no indices found for %s expression [total %d in the cluster]", args[0], len(indexNames))
				return
			} else {
				log.Info().Msgf(
					"found %d %s for %s expression:\n%s",
					len(filtered), fp.Ternary("index", "indices", len(filtered) == 1), args[0], strings.Join(filtered, "\n"))
				if flagutils.GetBoolFlag(cmd.Flags(), ConfirmFlag) ||
					prompts.IsOk(
						prompts.QuestionPrompt(
							fmt.Sprintf(
								"[context:%s]Are you sure you want to delete these indices?", client.Config.Current))) {
					for _, index := range filtered {
						log.Info().Msgf("deleting index '%s'", index)
						if err := client.DeleteIndex(index); err != nil {
							log.Fatal().Msgf("fail to delete index:%v", err)
						}
					}
				}
			}
		}
	},
}
