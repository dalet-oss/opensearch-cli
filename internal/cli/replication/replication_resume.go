package stats

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"sort"
)

const RawFlag = "raw"

var resumeReplicationCmd = &cobra.Command{
	Use:   "resume",
	Short: "resume replication",
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		replicationIndex := ""
		client := api.New(configutils.LoadConfig(appConfigFile))

	},
}

func init() {
	resumeReplicationCmd.PersistentFlags().Bool(RawFlag, false, "show raw api response")
}
