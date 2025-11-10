package consts

import (
	"os"
	"path"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

const (
	//ServiceName This is the service name in the keyring object
	ServiceName = "dalet-opensearch-cli"
	//CredSeparator is the separator used to separate the username and password in the credentials.
	CredSeparator = ":::"

	// DataDir is the default directory name for storing application-specific data and configuration files.
	DataDir = ".dalet"
	// Tooldir is the default directory name for storing application-specific data and configuration files.
	Tooldir = "oscli"
	// ConfigFile is the default name of the config file.
	ConfigFile = "config"

	// ConfigFlag used to override path to the cli config
	ConfigFlag = "config"
	// RawFlag signals to print raw API response from the OpenSearch cluster
	RawFlag = "raw"
	// VersionFlag print version of the cli and exit
	VersionFlag = "version"
	// VaultPasswordFlag Supply a vault password file to decrypt the vaulted credentials.
	VaultPasswordFlag = "vault-password"
	// DefaultRemoteClusterAlias is the default name alias used for identifying the remote OpenSearch cluster.
	DefaultRemoteClusterAlias = "pyramid-replication"
)

// bootstrapAndGet bootstraps the config dir and returns the path to it.
func bootstrapAndGet() string {
	dir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		log.Fatal().Msgf("❌Unable to get user home dir:%v", homeDirErr)
	}
	configPathDir := path.Join(dir, DataDir, Tooldir)
	err := os.MkdirAll(configPathDir, 0755)
	if err != nil {
		log.Fatal().Msgf("❌Unable to create config dir:%v", err)
	}
	return configPathDir
}

// DefaultConfig returns the path to the default config file.
func DefaultConfig() string {
	configPathDir := bootstrapAndGet()
	return path.Join(configPathDir, ConfigFile)
}
