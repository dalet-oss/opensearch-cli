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
)
import "github.com/dalet-oss/opensearch-cli/pkg/utils/logging"

var log = logging.Logger()

// OpensearchWrapper is a wrapper around the OpenSearch client that provides additional functionality.
// It provides methods for common operations on the cluster, such as retrieving cluster settings and listing plugins.
// It also provides a method for creating a new OpensearchWrapper instance using the provided cobra.Command for configuration and context.
// The wrapper is used by the CLI commands to perform operations on the cluster.
// The wrapper is also used by the unit tests to mock the OpenSearch client.
type OpensearchWrapper struct {
	Client *opensearch.Client
	Config appconfig.AppConfig
}

// requestContext initializes a context with a timeout defined by the ServerCallTimeout configuration.
// Returns a cancellable context and its cancel function.
func (api *OpensearchWrapper) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), api.Config.ServerCallTimeout())
}

// NewFromCmd creates a new OpensearchWrapper instance using the provided cobra.Command for configuration and context.
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

// New creates a new OpensearchWrapper instance using the provided appconfig.AppConfig and context.
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

// ClusterSettings retrieves and logs the current settings of the OpenSearch cluster, excluding default settings.
func (api *OpensearchWrapper) ClusterSettings() error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var rspData interface{}
	_, err := api.Client.Do(ctx, opensearchapi.ClusterGetSettingsReq{
		Params: opensearchapi.ClusterGetSettingsParams{IncludeDefaults: opensearch.ToPointer(false)},
	}, &rspData)
	if err != nil {
		return err
	}
	log.Info().Msgf("opensearch cluster settings:\n%s", printutils.MarshalJSONOrDie(rspData))
	return nil
}

// PluginsList retrieves and logs the list of installed plugins from the OpenSearch cluster.
func (api *OpensearchWrapper) PluginsList() ([]opensearchapi.CatPluginResp, error) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var rspData []opensearchapi.CatPluginResp
	_, err := api.Client.Do(ctx, opensearchapi.CatPluginsReq{Params: opensearchapi.CatPluginsParams{}}, &rspData)
	if err != nil {
		return nil, err

	}
	return rspData, nil
}
