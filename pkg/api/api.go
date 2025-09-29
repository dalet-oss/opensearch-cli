package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/appconfig"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/consts"
	configutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/config"
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/flagutils"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
	"context"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type ApiResponse map[string]interface{}

const LightOperationTimeout = 10 * time.Second

type OpensearchWrapper struct {
	Client *opensearch.Client
	Config appconfig.AppConfig
}

func NewFromCmd(cmd *cobra.Command) *OpensearchWrapper {
	return New(
		configutils.LoadConfig(flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag)),
		configutils.CreateApiContext(cmd),
	)
}

func New(c appconfig.AppConfig, ctx context.Context) *OpensearchWrapper {
	return &OpensearchWrapper{
		Client: GetOpenSearchClient(c, ctx),
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
