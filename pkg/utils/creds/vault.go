package creds

import (
	"errors"
	"github.com/dalet-oss/opensearch-cli/pkg/types"
	vault "github.com/sosedoff/ansible-vault-go"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// GetDataFromVaultString decrypts a Vault string and extracts the username and password based on specified keys.
func GetDataFromVaultString(vaultString, decryptPass, userKey, passKey string) (*types.Creds, error) {
	vaultYaml := make(map[string]interface{})
	result := new(types.Creds)
	decryptedString, decryptErr := vault.Decrypt(vaultString, decryptPass)
	if decryptErr != nil {
		return nil, errors.Join(errors.New("unable to decrypt vault string"), decryptErr)
	}
	if marshallErr := yaml.Unmarshal([]byte(decryptedString), vaultYaml); marshallErr != nil {
		return nil, marshallErr
	} else {
		if username, ok := vaultYaml[userKey].(string); ok {
			result.Username = username
		} else {
			return nil, errors.New("unable to find username in vault string")
		}
		if password, ok := vaultYaml[passKey].(string); ok {
			result.Password = password
		} else {
			return nil, errors.New("unable to find password in vault string")
		}
	}
	return result, nil
}

// GetDataFromVaultFile decrypts a Vault file and extracts the username and password based on specified keys.
func GetDataFromVaultFile(vaultFile, decryptPass, userKey, passKey string) (*types.Creds, error) {
	vaultYaml := make(map[string]interface{})
	result := new(types.Creds)
	if strings.HasPrefix(vaultFile, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err)
		}
		vaultFile = filepath.Join(homeDir, vaultFile[2:])
	}
	fileContent, decryptErr := vault.DecryptFile(vaultFile, decryptPass)
	if decryptErr != nil {
		return nil, errors.Join(errors.New("unable to decrypt vault file"), decryptErr)
	}
	if marshallErr := yaml.Unmarshal([]byte(fileContent), vaultYaml); marshallErr != nil {
		return nil, marshallErr
	} else {
		if username, ok := vaultYaml[userKey].(string); ok {
			result.Username = username
		} else {
			return nil, errors.New("unable to find username in vault file")
		}
		if password, ok := vaultYaml[passKey].(string); ok {
			result.Password = password
		} else {
			return nil, errors.New("unable to find password in vault file")
		}
	}
	return result, nil
}

// CreateVault encrypts a map of vaultData using the provided vaultPass and returns the encrypted vault string.
func CreateVault(vaultData map[string]string, vaultPass string) string {
	vaultYamlBytes, marshallErr := yaml.Marshal(vaultData)
	if marshallErr != nil {
		log.Fatal().Msgf("fail to marshal content:%v", marshallErr)
	}
	vaultString, encryptErr := vault.Encrypt(string(vaultYamlBytes), vaultPass)
	if encryptErr != nil {
		log.Fatal().Msgf("unable to encrypt vault content:%v", encryptErr)
	}
	return vaultString
}
