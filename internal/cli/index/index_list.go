package index

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"strings"
)

const FlagAll = "all"

var indexListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all indices.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		// method
		indices := api.New(configutils.LoadConfig(appConfigFile)).GetIndexList()
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

func init() {
	indexListCmd.Flags().Bool(FlagAll, false, "show all indices, including hidden ones(starting with '.').")
	indexCmd.AddCommand(indexListCmd)
}
