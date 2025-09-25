package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/appconfig"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
	"context"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"log"
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

// Generic methods for wrapper | not specific to plugin
func (a *OpensearchWrapper) ClusterSettings() {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var rspData interface{}
	_, err := a.Client.Do(ctx, opensearchapi.ClusterGetSettingsReq{
		Params: opensearchapi.ClusterGetSettingsParams{IncludeDefaults: opensearch.ToPointer(false)},
	}, &rspData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("opensearch cluster settings:\n%s", printutils.MarshalJSONOrDie(rspData))
}

func (a *OpensearchWrapper) PluginsList() []opensearchapi.CatPluginResp {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var rspData []opensearchapi.CatPluginResp
	_, err := a.Client.Do(ctx, opensearchapi.CatPluginsReq{Params: opensearchapi.CatPluginsParams{}}, &rspData)
	if err != nil {
		log.Fatalf("fail to get plugin list:%v", err)
	}
	//log.Printf("list of plugins:\n%s", printutils.MarshalJSONOrDie(rspData))
	return rspData
}
