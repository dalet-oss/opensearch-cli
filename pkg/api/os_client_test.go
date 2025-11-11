package api

import (
	"context"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/creds"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildOSConfig(t *testing.T) {
	type args struct {
		c   appconfig.AppConfig
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		ValidatorFunc func(*testing.T, opensearch.Config)
		WantErr       bool
	}{
		{
			name: "valid config",
			args: args{
				c:   *ConfigTContainer(opensearchContainer),
				ctx: contextWithPassword,
			},
			ValidatorFunc: nil,
			WantErr:       false,
		},
		{
			name: "invalid config",
			args: args{
				c:   appconfig.AppConfig{},
				ctx: nil,
			},
			ValidatorFunc: nil,
			WantErr:       true,
		},
		{
			name: "invalid config|no active context is set",
			args: args{
				c: appconfig.AppConfig{
					Current: "",
				},
				ctx: nil,
			},
			ValidatorFunc: nil,
			WantErr:       true,
		},
		{
			name: "invalid config|current set to non-existent context",
			args: args{
				c: appconfig.AppConfig{
					Current: "non-existent",
				},
				ctx: nil,
			},
			ValidatorFunc: nil,
			WantErr:       true,
		},
		{
			name: "invalid config|context is set, but cluster definition is missing or incorrect",
			args: args{
				c: appconfig.AppConfig{
					Current: cName,
					Contexts: []appconfig.ContextConfig{
						{
							Name:    cName,
							Cluster: cName,
							User:    cName,
						},
					},
				},
				ctx: nil,
			},
			ValidatorFunc: nil,
			WantErr:       true,
		},
		{
			name: "invalid config|context,cluster is ok, but user definition is missing or incorrect",
			args: args{
				c: appconfig.AppConfig{
					Current: cName,
					Contexts: []appconfig.ContextConfig{
						{
							Name:    cName,
							Cluster: cName,
							User:    cName,
						},
					},
					Clusters: []appconfig.ClusterConfig{
						{
							Name: cName,
							Params: appconfig.ClusterParams{
								Server: "http://localhost:9200",
								Tls:    false,
							},
						},
					},
				},
				ctx: nil,
			},
			WantErr: true,
		},
		{
			name: "invalid config|context,cluster,user is ok, but user creds are missing or incorrect",
			args: args{
				c: appconfig.AppConfig{
					Current: cName,
					Contexts: []appconfig.ContextConfig{
						{
							Name:    cName,
							Cluster: cName,
							User:    cName,
						},
					},
					Clusters: []appconfig.ClusterConfig{
						{
							Name: cName,
							Params: appconfig.ClusterParams{
								Server: "http://localhost:9200",
								Tls:    false,
							},
						},
					},
					Users: []appconfig.UserConfig{
						{
							Name: cName,
							User: appconfig.User{
								Vault: &appconfig.VaultConfig{
									VaultString: "",
									Username:    vaultUserKey,
									Password:    vaultPasswordKey,
								},
							},
						},
					},
				},
				ctx: nil,
			},
			WantErr: true,
		},
		{
			name: "invalid config|context,cluster,user is ok, but user creds are is incorrect",
			args: args{
				c: appconfig.AppConfig{
					Current: cName,
					Contexts: []appconfig.ContextConfig{
						{
							Name:    cName,
							Cluster: cName,
							User:    cName,
						},
					},
					Clusters: []appconfig.ClusterConfig{
						{
							Name: cName,
							Params: appconfig.ClusterParams{
								Server: "http://localhost:9200",
								Tls:    false,
							},
						},
					},
					Users: []appconfig.UserConfig{
						{
							Name: cName,
							User: appconfig.User{
								Vault: &appconfig.VaultConfig{
									VaultString: creds.CreateVault(map[string]string{
										vaultUserKey:     "superuser",
										vaultPasswordKey: "superpass",
									}, "veryStr0ngPassw0rd!"),
									Username: vaultUserKey,
									Password: vaultPasswordKey,
								},
							},
						},
					},
				},
				ctx: context.WithValue(context.TODO(), consts.VaultPasswordFlag, "wrong"),
			},
			WantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osConfig, err := BuildOSConfig(tt.args.c, tt.args.ctx)
			if tt.WantErr {
				assert.Error(t, err, "error is expected")
			} else {
				assert.NoError(t, err, "no error is expected")
			}
			if tt.ValidatorFunc != nil {
				tt.ValidatorFunc(t, osConfig)
			}
		})
	}
}

func TestGetOpenSearchClient(t *testing.T) {
	tests := []struct {
		name          string
		config        appconfig.AppConfig
		ctx           context.Context
		WantErr       bool
		ValidatorFunc func(*testing.T, *opensearch.Client)
	}{
		{
			name:    "perfect client",
			config:  *ConfigTContainer(opensearchContainer),
			ctx:     contextWithPassword,
			WantErr: false,
			ValidatorFunc: func(t *testing.T, cl *opensearch.Client) {
				assert.NotNil(t, cl)
			},
		},
		{
			name:    "nil client|empty config",
			config:  appconfig.AppConfig{},
			ctx:     nil,
			WantErr: true,
			ValidatorFunc: func(t *testing.T, cl *opensearch.Client) {
				assert.Nil(t, cl)
			},
		},
		{
			name:    "nil client|no creds",
			config:  *ConfigTContainer(opensearchContainer),
			ctx:     nil,
			WantErr: true,
			ValidatorFunc: func(t *testing.T, cl *opensearch.Client) {
				assert.Nil(t, cl)
			},
		},
		{
			name:    "nil client|wrong creds",
			config:  *ConfigTContainer(opensearchContainer),
			ctx:     context.WithValue(context.TODO(), consts.VaultPasswordFlag, "wrong"),
			WantErr: true,
			ValidatorFunc: func(t *testing.T, cl *opensearch.Client) {
				assert.Nil(t, cl)
			},
		},
		{
			name: "nil client|wrong server address",
			config: appconfig.AppConfig{
				Current: cName,
				Contexts: []appconfig.ContextConfig{
					{
						Name:    cName,
						Cluster: cName,
						User:    cName,
					},
				},
				Clusters: []appconfig.ClusterConfig{
					{
						Name: cName,
						Params: appconfig.ClusterParams{
							Server: "http://fake.local:9200",
							Tls:    false,
						},
					},
				},
				Users: []appconfig.UserConfig{
					{
						Name: cName,
						User: appconfig.User{
							Vault: &appconfig.VaultConfig{
								VaultString: creds.CreateVault(vaultData, vaultPassword),
								Username:    vaultUserKey,
								Password:    vaultPasswordKey,
							},
						},
					},
				},
			},
			ctx:     contextWithPassword,
			WantErr: true,
			ValidatorFunc: func(t *testing.T, cl *opensearch.Client) {
				assert.Nil(t, cl)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl, err := GetOpenSearchClient(tt.config, tt.ctx)
			if tt.WantErr {
				assert.Error(t, err, "error is expected")
			} else {
				assert.NoError(t, err, "no error is expected")
			}
			if tt.ValidatorFunc != nil {
				tt.ValidatorFunc(t, cl)
			}
		})
	}
}
