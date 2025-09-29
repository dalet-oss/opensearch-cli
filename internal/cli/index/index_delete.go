package index

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"slices"
	"sort"
)

const ConfirmFlag = "approve"

var indexDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "⚠️deletes index.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// method
		client := api.NewFromCmd(cmd)
		indexToDelete := ""
		registeredIndices := client.GetIndexList()
		indexNames := fp.Map(registeredIndices, func(info api.IndexInfo) string {
			return info.Index
		})
		sort.Strings(indexNames)
		if len(args) == 0 {
			// todo interactive mode
			indexToDelete = prompts.SelectivePrompt("Select index for removal", indexNames)
			return
		} else {
			// check if there's index -> delete
			indexToDelete = args[0]
			if idx := slices.IndexFunc(registeredIndices, func(info api.IndexInfo) bool {
				return info.Index == indexToDelete
			}); idx == -1 {
				log.Fatalf("❌index '%s' not found", indexToDelete)
			}

		}
		if prompts.IsOk(prompts.QuestionPrompt(fmt.Sprintf("Are you sure you want to delete index '%s'?", indexToDelete))) {
			client.DeleteIndex(indexToDelete)
		}
	},
}

func init() {
	indexListCmd.Flags().Bool(ConfirmFlag, false, "show all indices, including hidden ones(starting with '.').")
	indexCmd.AddCommand(indexDeleteCmd)
}
