package api

import (
	"context"
	"errors"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/network"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	aliases := []string{"tc-net"}
	net, err := network.New(ctx)
	if err != nil {
		log.Fatal().Err(err)
	}
	for _, containerName := range osContainers {
		osCtrx[containerName] = spinOpenSearch(aliases, net, containerName)
	}
	opensearchContainer = osCtrx[MainContainer]
	config = ConfigTContainer(opensearchContainer)
	defer func() {
		var errx error
		for kKey, vContainer := range osCtrx {
			if err := vContainer.Terminate(context.Background()); err != nil {
				errx = errors.Join(errx, err)
			}
			delete(osCtrx, kKey)
		}
		if errNet := net.Remove(ctx); errNet != nil {
			errx = errors.Join(errx, errNet)
		}
		if errx != nil {
			log.Fatal().Err(errx)
		}
	}()
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			c   appconfig.AppConfig
			ctx context.Context
		}
		wantErr bool
	}{
		{
			name: "auto-created OpenSearch client",
			args: struct {
				c   appconfig.AppConfig
				ctx context.Context
			}{
				c:   *config,
				ctx: contextWithPassword,
			},
			wantErr: false,
		},
	}
	assert.NotNil(t, opensearchContainer)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, clientErr := New(tt.args.c, tt.args.ctx)
			assert.NotNil(t, client)
			if tt.wantErr {
				assert.Error(t, clientErr, "error creating client is expected")
			} else {
				assert.NoError(t, clientErr, "no client creation error is expected")
			}
		})
	}
}

func TestOpensearchWrapper_ClusterSettings(t *testing.T) {
	tests := []OSSingleContainerTest{
		{
			Name:    "get settings of the normal cluster",
			Wrapper: testWrapper(),
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if !tt.WantErr {
				assert.NoError(t, tt.Wrapper.ClusterSettings(), "expected to get cluster settings")
			} else {
				assert.Error(t, tt.Wrapper.ClusterSettings(), "expected to get error")
			}
		})
	}
}

func TestOpensearchWrapper_PluginsList(t *testing.T) {
	tests := []OSSingleContainerTest{
		{
			Name:    "get list of plugins",
			Wrapper: testWrapper(),
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			plugins, queryPluginErr := tt.Wrapper.PluginsList()
			if !tt.WantErr {
				assert.NoError(t, queryPluginErr, "expected to get list of plugins")
				assert.NotEmpty(t, plugins, "expected some plugins to be available")
			} else {
				assert.Error(t, queryPluginErr, "expected to get error")
			}
		})
	}
}
