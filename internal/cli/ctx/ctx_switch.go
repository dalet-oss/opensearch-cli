/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package ctx

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	configutils "github.com/dalet-oss/opensearch-cli/pkg/utils/config"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/prompts"
	"github.com/spf13/cobra"
)

// ctxSwitchCmd represents the switch command
var ctxSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switches to a context.",
	Long: `
Switch to active context to the one chosen by the user.
If context name is not provided, it will prompt for the context name(from the list of contexts).
`,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile, _ := cmd.Flags().GetString(consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		if len(args) > 0 {
			if config.HasContext(appconfig.ContextConfig{Name: args[0]}) {
				config.Current = args[0]
			} else {
				fmt.Printf("requested context '%s' is not found.\n", args[0])
				return
			}
		} else {
			config.Current = prompts.SelectivePrompt("Select context", config.GetContextList())
		}
		if !configutils.SaveConfig(appConfigFile, config) {
			fmt.Println("❌Failed to save config file. Please try again.")
		}
	},
}
