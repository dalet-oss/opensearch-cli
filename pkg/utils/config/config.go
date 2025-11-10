package config

import (
	"context"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/logging"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var log = logging.Logger()

const ConfigFilePerm = 0644

// configBytes returns the config as a byte array.
// The config is marshalled to YAML format.
// If the marshaling fails, the program is terminated with an error.
func configBytes(config appconfig.AppConfig) []byte {
	marshal, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal().Msgf("unable to marshal config:%v", err)
	}
	return marshal
}

// writeFile writes data to a specified file at the given path with predefined file permissions. Logs a fatal error on failure.
func writeFile(path string, data []byte) {
	writeErr := os.WriteFile(path, data, ConfigFilePerm)
	if writeErr != nil {
		log.Fatal().Msgf("unable to write config file:%v", writeErr)
	}
}

// writeFileResult writes data to a specified file at the given path with predefined file permissions. Logs a fatal error on failure.
// Returns true if the file was written successfully, false otherwise.
func writeFileResult(path string, data []byte) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			log.Fatal().Msgf("unable to create config file:%v", err)
		}
	}
	writeErr := os.WriteFile(path, data, ConfigFilePerm)
	if writeErr != nil {
		log.Info().Msgf("unable to write config file:%v", writeErr)
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
			log.Fatal().Msgf("unable to create config file:%v", err)
		}
		if example {
			writeFile(defaultConfig, configBytes(appconfig.Example()))
		} else {
			writeFile(defaultConfig, configBytes(appconfig.AppConfig{ApiVersion: appconfig.ApiVersionV1}))
		}

	}
}

// LoadConfig loads the application configuration from the specified path or defaults to the predefined config path.
// If the file does not exist, it initializes a new configuration file if using the default path.
// If an error occurs while reading or unmarshalling the file, the function logs a fatal error.
// Returns the loaded AppConfig structure.
func LoadConfig(path string) appconfig.AppConfig {
	configPath := consts.DefaultConfig()
	if len(path) > 0 {
		configPath = path
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) && configPath == consts.DefaultConfig() {
		Init(false)
	} else if os.IsNotExist(err) {
		log.Fatal().Msgf("config file not found:%v", err)
	}
	if fileContent, err := os.ReadFile(configPath); err != nil {
		log.Fatal().Msgf("unable to read config file:%v", err)
	} else {
		var config appconfig.AppConfig
		marshalErr := yaml.Unmarshal(fileContent, &config)
		if marshalErr != nil {
			log.Fatal().Msgf("unable to unmarshal config file:%v", marshalErr)
		}
		return config
	}
	return appconfig.AppConfig{}
}

// SaveConfig saves the application configuration to the specified path or defaults to the predefined config path.
// If the file does not exist, it initializes a new configuration file if using the default path.
// If an error occurs while writing the file, the function logs a fatal error.
// Returns true if the file was written successfully, false otherwise.
func SaveConfig(path string, config appconfig.AppConfig) bool {
	configPath := consts.DefaultConfig()
	if len(path) > 0 {
		configPath = path
	}
	return writeFileResult(configPath, configBytes(config))
}

// CreateApiContext returns a context enriched with the vault password flag value if it is set; otherwise, the base context.
func CreateApiContext(cmd *cobra.Command) context.Context {
	if cmd.Flags().Lookup(consts.VaultPasswordFlag).Changed {
		return context.WithValue(cmd.Context(), consts.VaultPasswordFlag, flagutils.GetStringFlag(cmd.Flags(), consts.VaultPasswordFlag))
	}
	return cmd.Context()
}
