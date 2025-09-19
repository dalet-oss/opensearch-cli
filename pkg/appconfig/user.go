package appconfig

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/types"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/creds"
	"fmt"
)

// UserConfig represents the configuration for managing a single user.
type UserConfig struct {
	// Name holds the username.
	Name string `yaml:"name"`
	// User entity.
	User User `yaml:"user"`
}

// User represents the user entity.
type User struct {
	// Token represents a security token used for authentication, serialized as "token" in YAML, and omitted if empty.
	Token string `yaml:"token,omitempty"`
	// Exec represents a command to execute for authentication, serialized as "exec" in YAML, and omitted if empty.
	Exec CredentialsExec `yaml:"exec,omitempty"`
}

type CredentialsExec struct {
	Cmd  string        `yaml:"cmd"`
	Args []interface{} `json:"args"`
}

func (u *UserConfig) GetUserCredentials() (types.Creds, error) {
	if u.User.Token != "" {
		return creds.PullFromKeyring(u.User.Token), nil
	}
	return types.Creds{}, fmt.Errorf("not implemented")
}
