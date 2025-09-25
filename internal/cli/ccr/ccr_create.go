package ccr

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	SettingsTypeFlag       = "type"
	SettingsModeFlag       = "remote-mode"
	SettingsRemoteNameFlag = "remote-name"
	SettingsRemoteAddrFlag = "remote-addr"
)

var ccrCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ccr object in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		api.New(
			configutils.LoadConfig(
				flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag))).
			ConfigureRemoteCluster(prepareOpts(cmd.Flags()), flagutils.GetBoolFlag(cmd.Flags(), consts.RawFlag))
	},
}

func init() {
	ccrCreateCmd.PersistentFlags().String(SettingsTypeFlag, "", "type of the settings [transient,persistent,default]")
	ccrCreateCmd.PersistentFlags().String(SettingsModeFlag, "", "remote mode:[proxy,sniff]")
	ccrCreateCmd.PersistentFlags().String(SettingsRemoteNameFlag, consts.DefaultRemoteClusterAlias, "remote name alias")
	ccrCreateCmd.PersistentFlags().String(SettingsRemoteAddrFlag, "", "address of the remote cluster http(s)://<host>[:port]")
}

func prepareOpts(flags *pflag.FlagSet) api.CCRCreateOpts {
	opts := api.CCRCreateOpts{
		Type:       flagutils.GetStringFlagInSet(flags, SettingsTypeFlag, []string{"transient", "persistent", "default", ""}),
		Mode:       flagutils.GetStringFlagInSet(flags, SettingsModeFlag, []string{"proxy", ""}),
		RemoteName: flagutils.GetNotEmptyStringFlag(flags, SettingsRemoteNameFlag),
		RemoteAddr: flagutils.GetNotEmptyStringFlag(flags, SettingsRemoteAddrFlag),
	}
	return opts
}
