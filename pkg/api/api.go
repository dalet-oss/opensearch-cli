package api

import (
	"context"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	configutils "github.com/dalet-oss/opensearch-cli/pkg/utils/config"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/flagutils"
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/spf13/cobra"
	"time"
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

type ApiResponse map[string]interface{}

const LightOperationTimeout = 10 * time.Second

type OpensearchWrapper struct {
	Client *opensearch.Client
	Config appconfig.AppConfig
}

func (a *OpensearchWrapper) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), a.Config.ServerCallTimeout())
}

func NewFromCmd(cmd *cobra.Command) *OpensearchWrapper {
	wrapper, err := New(
		configutils.LoadConfig(flagutils.GetStringFlag(cmd.Flags(), consts.ConfigFlag)),
		configutils.CreateApiContext(cmd),
	)
	if err != nil {
		log.Warn().Msg("unable to create client, check your config file.")
		log.Fatal().Err(err)
	}
	return wrapper
}

func New(c appconfig.AppConfig, ctx context.Context) (*OpensearchWrapper, error) {
	client, err := GetOpenSearchClient(c, ctx)
	if err != nil {
		return nil, err
	}
	return &OpensearchWrapper{
		Client: client,
		Config: c,
	}, nil
}

// Generic methods for wrapper | not specific to plugin
func (a *OpensearchWrapper) ClusterSettings() error {
	ctx, cancelFunc := a.requestContext()
	defer cancelFunc()
	var rspData interface{}
	_, err := a.Client.Do(ctx, opensearchapi.ClusterGetSettingsReq{
		Params: opensearchapi.ClusterGetSettingsParams{IncludeDefaults: opensearch.ToPointer(false)},
	}, &rspData)
	if err != nil {
		return err
	}
	log.Info().Msgf("opensearch cluster settings:\n%s", printutils.MarshalJSONOrDie(rspData))
	return nil
}

func (a *OpensearchWrapper) PluginsList() ([]opensearchapi.CatPluginResp, error) {
	ctx, cancelFunc := a.requestContext()
	defer cancelFunc()
	var rspData []opensearchapi.CatPluginResp
	_, err := a.Client.Do(ctx, opensearchapi.CatPluginsReq{Params: opensearchapi.CatPluginsParams{}}, &rspData)
	if err != nil {
		return nil, err

	}
	return rspData, nil
}
