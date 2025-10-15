package userconfig

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/appconfig"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/creds"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/generic"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/logging"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/prompts"
	"fmt"
)

var log = logging.Logger()

const (
	VaultBackend   = "vault"
	KeyringBackend = "keyring"
)

func CreateUserEntry(conf appconfig.AppConfig, cluster appconfig.ClusterConfig) appconfig.UserConfig {
	user := appconfig.UserConfig{}
	user.Name = prompts.ValidatedPrompt("User entry name", func(input string) error {
		if len(input) == 0 {
			return fmt.Errorf("user entry name is required")
		}
		if conf.HasUser(appconfig.UserConfig{Name: input}) {
			return fmt.Errorf("user entry name '%s' already exists", input)
		}
		return nil
	})
	switch prompts.SelectivePrompt("Choose credential backend", []string{VaultBackend, KeyringBackend}) {
	case VaultBackend:
		user.User = appconfig.User{
			Vault: getDataForVault(),
		}
	case KeyringBackend:
		username, pass := getDataForKeyring()
		if user.Name == "" {
			user.Name = fmt.Sprintf("%s@%s", username, cluster.Name)
		}
		user.User = appconfig.User{
			Token: creds.PushToKeyring(username, pass),
		}
	}
	return user
}

func getDataForKeyring() (string, string) {
	username := ""
	if userName := prompts.ValidatedPrompt("Username", func(input string) error {
		if len(input) == 0 {
			return fmt.Errorf("username is required")
		}
		return nil
	}); len(userName) == 0 {
		log.Fatal().Msg("Username is required")
	} else {
		username = userName
	}
	return username, prompts.SecretPrompt("Password")
}

func getDataForVault() *appconfig.VaultConfig {
	const bindAction = "bind"
	const createAction = "create"
	actions := []string{bindAction, createAction}
	switch prompts.SelectivePrompt("Create vault or bind to existing", actions) {
	case bindAction:
		return &appconfig.VaultConfig{
			File:        prompts.ValidatedPrompt("Vault file location", prompts.NotEmptyString),
			VaultString: "",
			Username:    prompts.ValidatedPrompt("Vault username key", prompts.NotEmptyString),
			Password:    prompts.ValidatedPrompt("Vault password key", prompts.NotEmptyString),
		}
	case createAction:
		userKey := prompts.ValidatedPrompt("Vault username key", prompts.NotEmptyString)
		username := prompts.ValidatedPrompt("username", prompts.NotEmptyString)
		passKey := prompts.ValidatedPrompt("Vault password key", prompts.NotEmptyString)
		password := prompts.SecretPrompt("password")
		vaultConfig := &appconfig.VaultConfig{
			Username: userKey,
			Password: passKey,
		}
		vaultCreds := make(map[string]string)
		vaultCreds[vaultConfig.Username] = username
		vaultCreds[vaultConfig.Password] = password
		const embedVault = "embed"
		const saveVault = "save"
		vaultString := creds.CreateVault(vaultCreds, prompts.SecretPrompt("Set Vault password"))
		switch prompts.SelectivePrompt("How do you want to save vault credentials?", []string{embedVault, saveVault}) {
		case embedVault:
			vaultConfig.VaultString = vaultString
		case saveVault:
			vaultLocation := prompts.ValidatedPrompt("Choose vault location", prompts.NotEmptyString)
			generic.SaveFileOrDie(vaultLocation, vaultString)
			vaultConfig.File = vaultLocation
		}
		return vaultConfig
	}
	return nil
}
