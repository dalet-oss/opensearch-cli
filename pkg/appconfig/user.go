package appconfig

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/types"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/creds"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"context"
	"fmt"
)

// UserConfig represents the configuration for managing a single user.
type UserConfig struct {
	// Name holds the username.
	Name string `yaml:"name"`
	// User entity.
	User User `yaml:"user"`

	decryptedCreds *types.Creds `yaml:"-"`
}

// User represents the user entity.
type User struct {
	// Token represents a security token used for authentication, serialized as "token" in YAML, and omitted if empty.
	Token string `yaml:"token,omitempty"`
	// Exec represents a command to execute for authentication, serialized as "exec" in YAML, and omitted if empty.
	Exec CredentialsExec `yaml:"exec,omitempty"`
	// Vault represents a Vault authentication configuration
	Vault *VaultConfig `yaml:"vault,omitempty"`
}

type VaultConfig struct {
	// Path to the ansible-vault encrypted file with credentials.
	File string `yaml:"file"`
	// VaultString vault file content [optional, if file is not provided]
	VaultString string `yaml:"vault_string,omitempty"`
	// Username the key contains the username in the file
	Username string `yaml:"username"`
	// Password the key contains the password in the file
	Password string `yaml:"password"`
}

type CredentialsExec struct {
	Cmd  string        `yaml:"cmd"`
	Args []interface{} `json:"args"`
}

func (u *UserConfig) GetUserCredentials(ctx context.Context) (*types.Creds, error) {
	if u.decryptedCreds != nil {
		return u.decryptedCreds, nil
	}
	if u.User.Token != "" {
		u.decryptedCreds = creds.PullFromKeyring(u.User.Token)
		return u.decryptedCreds, nil
	}
	if u.User.Vault != nil {
		if u.User.Vault.File == "" && u.User.Vault.VaultString == "" {
			return nil, fmt.Errorf("vault file or vault string must be provided")
		}
		if u.User.Vault.Username == "" || u.User.Vault.Password == "" {
			return nil, fmt.Errorf("vault username and password keys must be provided")
		}
		vaultPassword := ""
		if v := ctx.Value(consts.VaultPasswordFlag); v != nil {
			vaultPassword = v.(string)
		} else {
			vaultPassword = prompts.SecretPrompt("Vault password:")
		}
		if u.User.Vault.File != "" {
			if vaultCreds, err := creds.GetDataFromVaultFile(
				u.User.Vault.File,
				vaultPassword,
				u.User.Vault.Username,
				u.User.Vault.Password); err != nil {
				return nil, err
			} else {
				u.decryptedCreds = vaultCreds
				return u.decryptedCreds, nil
			}
		} else {
			vaultCreds, err := creds.GetDataFromVaultString(
				u.User.Vault.VaultString,
				vaultPassword,
				u.User.Vault.Username,
				u.User.Vault.Password)
			if err != nil {
				return nil, err
			} else {
				u.decryptedCreds = vaultCreds
				return u.decryptedCreds, nil
			}
		}
	}
	return nil, fmt.Errorf("not implemented")
}
