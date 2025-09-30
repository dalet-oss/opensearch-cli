package consts

import (
	"log"
	"os"
	"path"
)

const (
	//ServiceName This is the service name in the keyring object
	ServiceName   = "dalet-opensearch-cli"
	CredSeparator = ":::"
	DataDir       = ".dalet"
	Tooldir       = "oscli"
	ConfigFile    = "config"

	// ConfigFlag used to override path to the cli config
	ConfigFlag = "config"
	// RawFlag signals to print raw API response from the OpenSearch cluster
	RawFlag = "raw"
	// VersionFlag print version of the cli and exit
	VersionFlag = "version"
	// VaultPasswordFlag Supply a vault password file to decrypt the vaulted credentials.
	VaultPasswordFlag         = "vault-password"
	DefaultRemoteClusterAlias = "pyramid-replication"
)

// bootstrapAndGet bootstraps the config dir and returns the path to it.
func bootstrapAndGet() string {
	dir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		log.Fatalf("❌Unable to get user home dir:%v", homeDirErr)
	}
	configPathDir := path.Join(dir, DataDir, Tooldir)
	err := os.MkdirAll(configPathDir, 0755)
	if err != nil {
		log.Fatalf("❌Unable to create config dir:%v", err)
	}
	return configPathDir
}

// DefaultConfig returns the path to the default config file.
func DefaultConfig() string {
	configPathDir := bootstrapAndGet()
	return path.Join(configPathDir, ConfigFile)
}
