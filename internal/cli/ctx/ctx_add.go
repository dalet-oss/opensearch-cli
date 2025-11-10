/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package ctx

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	configutils "github.com/dalet-oss/opensearch-cli/pkg/utils/config"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/creds"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/prompts"
	"github.com/dalet-oss/opensearch-cli/pkg/ux/userconfig"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// TODO:
// Flags:
// -c, --cluster string   cluster url
// -u, --user string      user name
// -p, --password string  password
//
// Example:
// opensearch-cli ctx add -c https://localhost:9200 -u admin -p admin
//
// If flags are not passed, it will prompt for the values
// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "a"},
	Short:   "Creates a new context",
	Long:    `Creates a new context in the config file using interactive mode`,
	Run: func(cmd *cobra.Command, args []string) {
		appConfigFile := flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag)
		config := configutils.LoadConfig(appConfigFile)
		newCluster := CreateClusterEntry(config)
		user := userconfig.CreateUserEntry(config, newCluster)
		ctx := CreateContextEntry(config, newCluster, user)
		config.Push(newCluster, user, ctx)
		if prompts.IsOk(prompts.QuestionPrompt("Do you want to switch to the created context?")) {
			config.Current = ctx.Name
		}
		if !configutils.SaveConfig(appConfigFile, config) {
			creds.DeleteFromKeyring(user.User.Token)
		}
	},
}

// CreateClusterEntry creates a new cluster entry
func CreateClusterEntry(conf appconfig.AppConfig) appconfig.ClusterConfig {
	clusterConfig := appconfig.ClusterConfig{}
	if clusterName := prompts.ValidatedPrompt("(optional)Cluster name", func(input string) error {
		if len(input) == 0 {
			return nil
		} else {
			if conf.HasCluster(appconfig.ClusterConfig{Name: input}) {
				return fmt.Errorf("cluster name '%s' already exists", input)
			}
		}
		return nil
	}); len(clusterName) > 0 {
		clusterConfig.Name = clusterName
	}
	if clusterUrl := prompts.SimplePrompt("Cluster url"); len(clusterUrl) == 0 {
		log.Fatal().Msg("Cluster url is required")
	} else {
		clusterConfig.Params = appconfig.ClusterParams{
			Server: clusterUrl,
		}
	}
	if enableTls := prompts.ValidatedPrompt("Skip TLS verify(t|true|f|false)", func(input string) error {
		_, err := strconv.ParseBool(input)
		return err
	}); len(enableTls) == 0 {
		clusterConfig.Params.Tls = false
	} else {
		v, _ := strconv.ParseBool(enableTls)
		clusterConfig.Params.Tls = v
	}

	if len(clusterConfig.Name) == 0 {
		clusterConfig.Name = strings.ReplaceAll(strings.ReplaceAll(clusterConfig.Params.Server, "://", "::"), ":", "::")
	}
	return clusterConfig
}

// CreateContextEntry creates a new context entry
func CreateContextEntry(conf appconfig.AppConfig, cluster appconfig.ClusterConfig, user appconfig.UserConfig) appconfig.ContextConfig {
	newContext := appconfig.ContextConfig{}
	if auto := prompts.QuestionPrompt("Create context automatically?"); prompts.IsOk(auto) {
		newContext.Name = fmt.Sprintf("%s@%s", user.Name, cluster.Name)
		newContext.Cluster = cluster.Name
		newContext.User = user.Name
	} else {
		if newContexName := prompts.ValidatedPrompt("(optional)Context name:", func(input string) error {
			if len(input) == 0 {
				return nil
			}
			if conf.HasContext(appconfig.ContextConfig{Name: input}) {
				return fmt.Errorf("context name '%s' already exists", input)
			}
			return nil
		}); len(newContexName) == 0 {
			newContext.Name = fmt.Sprintf("%s@%s", user.Name, cluster.Name)
		} else {
			newContext.Name = newContexName
			newContext.Cluster = cluster.Name
			newContext.User = user.Name
		}
	}
	return newContext
}
