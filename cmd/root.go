package cmd

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/internal/cli/autofollow"
	"github.com/dalet-oss/opensearch-cli/internal/cli/ccr"
	"github.com/dalet-oss/opensearch-cli/internal/cli/ctx"
	"github.com/dalet-oss/opensearch-cli/internal/cli/index"
	"github.com/dalet-oss/opensearch-cli/internal/cli/replication"
	"github.com/dalet-oss/opensearch-cli/internal/cli/stats"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/config"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/logging"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var log = logging.Logger()
var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opensearch-cli",
	Short: "manage OpenSearch clusters and their indices.",
	Example: fmt.Sprintf(`
Some commands supports the operations with the wildcard(due to it's not available option in the OpenSearch server).
	%s
	`, gu.WildHelp),
	Run: func(cmd *cobra.Command, args []string) {
		// init things
		if flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag) == "" {
			config.Init(false)
		}
		if flagutils.GetBoolFlag(cmd.Flags(), consts.VersionFlag) {
			cmd.Printf("Version: %s\n\n", Version)
		} else {
			if err := cmd.Help(); err != nil {
				log.Err(err).Msg("failed to show help")
			}
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
		NewGendocCmd(rootCmd),
		ctx.NewCtxCmd(),
		index.NewIndexCmd(),
		stats.NewStatsCmd(),
		ccr.NewCCRCmd(),
		autofollow.NewAutofollowCmd(),
		replication.NewReplicationCmd(),
	)
}

// NewGendocCmd - generates simple markdown doc
func NewGendocCmd(rootCmd *cobra.Command) *cobra.Command {
	var docDir string
	var docCmd = &cobra.Command{
		Use:    "gendoc",
		Short:  "Generate Markdown documentation for the app",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat("docs"); os.IsExist(err) {
				_ = os.RemoveAll("docs")
			}

			err := os.MkdirAll("docs", os.ModePerm)
			if err != nil {

				log.Err(err).Msg("Failed to create docs directory")

			}
			if docDir == "" {
				docDir = "docs"
			}
			toMd(rootCmd, docDir)
		},
	}
	docCmd.Flags().StringVar(&docDir, "doc-output-dir", "", "the directory to output the docs to")
	return docCmd
}

func toMd(child *cobra.Command, outputDir string) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Err(err).Msg("Failed to create docs directory")
		}
	}
	filePath := path.Join(outputDir, child.Name()+".md")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal().Err(err)
	}

	capturedFunc := func(s string) string {
		isChildCmd := strings.Contains(s, child.Name())
		//captured here to capture the current variables needed to output the path in the callback
		parts := strings.Split(s, "_")
		partsLen := len(parts)
		n := parts[partsLen-1]
		n = n[:strings.Index(n, ".")]
		if isChildCmd {
			return path.Join(n, n+".md")
		} else {
			return path.Join("..", n+".md")
		}
	}
	err = doc.GenMarkdownCustom(child, file, capturedFunc)
	if err != nil {
		log.Fatal().Err(err)
	}

	_ = file.Close()

	for _, subCmd := range child.Commands() {
		if !subCmd.Hidden {
			toMd(subCmd, path.Join(outputDir, subCmd.Name()))
		}
	}
}
