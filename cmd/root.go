/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package cmd

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/autofollow"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/ccr"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/ctx"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/index"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/replication"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/stats"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opensearch-cli",
	Short: "manage OpenSearch clusters and their indices.",
	Run: func(cmd *cobra.Command, args []string) {
		// init things
		if flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag) == "" {
			config.Init(false)
		}
		if flagutils.GetBoolFlag(cmd.Flags(), consts.VersionFlag) {
			cmd.Printf("Version: %s\n\n", Version)
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// global flags
	rootCmd.PersistentFlags().Bool(consts.VersionFlag, false, "show version")
	rootCmd.PersistentFlags().String(consts.ConfigFlag, "", "config file (default is $HOME/.dalet/oscli/config)")
	rootCmd.PersistentFlags().String(consts.VaultPasswordFlag, "", "vault password for decrypting vault credentials")
	rootCmd.PersistentFlags().Bool(consts.RawFlag, false, "show raw api response")
	// subcommands
	rootCmd.AddCommand(
		ctx.NewCtxCmd(),
		index.NewIndexCmd(),
		stats.NewStatsCmd(),
		ccr.NewCCRCmd(),
		cli.NewClusterCmd(),
		autofollow.NewAutofollowCmd(),
		replication.NewReplicationCmd(),
	)
}
