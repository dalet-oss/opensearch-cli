package index

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/spf13/cobra"
)

var indexCreateCmd = &cobra.Command{
	Use:   "create",
	Short: " creates index",
	Long:  `Create index in the OpenSearch cluster`,
	Example: `
opensearch-cli index create <Index Name>
`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewFromCmd(cmd)
		if len(args) == 0 || args[0] == "" {
			log.Fatal().Msg("index name is required")
		} else if gu.ContainsWildcard(args[0]) {
			log.Fatal().Msg("wildcard is not allowed")
		} else {
			if err := client.CreateIndex(args[0]); err != nil {
				log.Fatal().Err(err)
			}
		}
	},
}
