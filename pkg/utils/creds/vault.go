package creds

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/types"
	"errors"
	vault "github.com/sosedoff/ansible-vault-go"
	"gopkg.in/yaml.v3"
)

func GetDataFromVaultString(vaultString, decryptPass, userKey, passKey string) (*types.Creds, error) {
	var vaultYaml map[string]interface{}
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

func GetDataFromVaultFile(vaultFile, decryptPass, userKey, passKey string) (*types.Creds, error) {
	var vaultYaml map[string]interface{}
	result := new(types.Creds)
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
