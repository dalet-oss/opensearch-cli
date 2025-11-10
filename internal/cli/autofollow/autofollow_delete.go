package autofollow

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var autofollowDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete autofollow rule from the cluster",
	Example: `opensearch-cli autofollow delete <RULE NAME> -l leader`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
			return
		}
		client := api.NewFromCmd(cmd)
		if err := client.
			DeleteAutofollow(autofollowDeleteOpts(cmd.Flags(), args[0]), flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag)); err != nil {
			log.Fatal().Msgf("failed to delete autofollow rule:%v", err)
		}
	},
}

func autofollowDeleteOpts(flags *pflag.FlagSet, name string) replication.DeleteAutofollowReq {
	opts := replication.DeleteAutofollowReq{}
	opts.Body = replication.DeleteAutofollowBody{
		Name:        name,
		LeaderAlias: flagutils.GetNotEmptyStringFlag(flags, LeaderAliasFlag),
	}
	return opts
}

func init() {
	autofollowDeleteCmd.PersistentFlags().String(LeaderAliasFlag, "", "leader alias")
}
