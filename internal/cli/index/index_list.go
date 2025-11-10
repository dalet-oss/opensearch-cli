package index

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"strings"
)

const FlagAll = "all"

var indexListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list all indices.",
	Long:    `show all indices in the OpenSearch cluster`,
	Example: `opensearch-cli index [list|ls]`,
	Run: func(cmd *cobra.Command, args []string) {
		indices, indexListErr := api.NewFromCmd(cmd).GetIndexList()
		if indexListErr != nil {
			log.Fatal().Err(indexListErr)
		}
		if v, _ := cmd.Flags().GetBool(FlagAll); !v {
			indices = slices.DeleteFunc(indices, func(info api.IndexInfo) bool {
				return strings.HasPrefix(info.Index, ".")
			})
		}
		slices.SortFunc(indices, func(a, b api.IndexInfo) int {
			return strings.Compare(a.Index, b.Index)
		})
		for _, index := range indices {
			fmt.Println(index)
		}
	},
}
