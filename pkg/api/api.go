package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/appconfig"
	"github.com/opensearch-project/opensearch-go/v4"
	"time"
)

type ApiResponse map[string]interface{}

const LightOperationTimeout = 10 * time.Second

type OpensearchWrapper struct {
	Client *opensearch.Client
	Config appconfig.AppConfig
}

func New(c appconfig.AppConfig) *OpensearchWrapper {
	return &OpensearchWrapper{
		Client: GetOpenSearchClient(c),
		Config: c,
	}
}
