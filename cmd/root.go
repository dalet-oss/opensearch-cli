/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package cmd

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/ctx"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/index"
	"bitbucket.org/ooyalaflex/opensearch-cli/internal/cli/stats"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opensearch-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.dalet/oscli/config)")
	rootCmd.AddCommand(ctx.NewCtxCmd())
	rootCmd.AddCommand(index.NewIndexCmd())
	rootCmd.AddCommand(stats.NewStatsCmd())
	rootCmd.AddCommand(cli.NewClusterCmd())
}
