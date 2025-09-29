package config

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/appconfig"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	"context"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const ConfigFilePerm = 0644

func configBytes(config appconfig.AppConfig) []byte {
	marshal, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("unable to marshal config:%v", err)
	}
	return marshal
}
func writeFile(path string, data []byte) {
	writeErr := os.WriteFile(path, data, ConfigFilePerm)
	if writeErr != nil {
		log.Fatalf("unable to write config file:%v", writeErr)
	}
}
func writeFileResult(path string, data []byte) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			log.Fatalf("unable to create config file:%v", err)
		}
	}
	writeErr := os.WriteFile(path, data, ConfigFilePerm)
	if writeErr != nil {
		log.Printf("unable to write config file:%v", writeErr)
		return false
	}
	return true
}

// Init initializes the config file.
// example - if true, the example config file is created.
func Init(example bool) {
	defaultConfig := consts.DefaultConfig()
	if _, err := os.Stat(defaultConfig); os.IsNotExist(err) {
		_, err := os.Create(defaultConfig)
		if err != nil {
			log.Fatalf("unable to create config file:%v", err)
		}
		if example {
			writeFile(defaultConfig, configBytes(appconfig.Example()))
		} else {
			writeFile(defaultConfig, configBytes(appconfig.AppConfig{ApiVersion: appconfig.ApiVersionV1}))
		}

	}
}

func LoadConfig(path string) appconfig.AppConfig {
	configPath := consts.DefaultConfig()
	if len(path) > 0 {
		configPath = path
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) && configPath == consts.DefaultConfig() {
		Init(false)
	} else if os.IsNotExist(err) {
		log.Fatalf("config file not found:%v", err)
	}
	if fileContent, err := os.ReadFile(configPath); err != nil {
		log.Fatalf("unable to read config file:%v", err)
	} else {
		var config appconfig.AppConfig
		marshalErr := yaml.Unmarshal(fileContent, &config)
		if marshalErr != nil {
			log.Fatalf("unable to unmarshal config file:%v", marshalErr)
		}
		return config
	}
	return appconfig.AppConfig{}
}

func SaveConfig(path string, config appconfig.AppConfig) bool {
	configPath := consts.DefaultConfig()
	if len(path) > 0 {
		configPath = path
	}
	return writeFileResult(configPath, configBytes(config))
}

func CreateApiContext(cmd *cobra.Command) context.Context {
	if cmd.Flags().Lookup(consts.VaultPasswordFlag).Changed {
		return context.WithValue(cmd.Context(), consts.VaultPasswordFlag, flagutils.GetStringFlag(cmd.Flags(), consts.VaultPasswordFlag))
	}
	return cmd.Context()
}
