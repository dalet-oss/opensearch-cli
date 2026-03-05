package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"net/http"
	"net/url"
)

// BuildOSConfig constructs and returns an OpenSearch configuration based on the given app configuration.
// It retrieves the active context, cluster, and user details, and incorporates user credentials into the configuration.
// The function panics if required elements like contexts, clusters, or users are missing or invalid.
func BuildOSConfig(c appconfig.AppConfig, ctx context.Context) (opensearch.Config, error) {
	ccfg := c.GetActiveContext()
	if ccfg == nil {
		return opensearch.Config{},
			fmt.Errorf("context config is not found, active context is set to '%s'", c.Current)
	}
	osConnection := ccfg.GetCluster(c)
	if osConnection == nil {
		return opensearch.Config{},
			fmt.Errorf("cluster definition '%s' is not found", ccfg.Cluster)
	}
	osUser := ccfg.GetUser(c)
	if osUser == nil {
		return opensearch.Config{},
			fmt.Errorf("user definition '%s' is not found", ccfg.User)
	}
	userCreds, err := osUser.GetUserCredentials(ctx)
	if err != nil {
		log.Warn().Msg("Unable to get user credentials, check your config file.")
		return opensearch.Config{}, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: osConnection.Params.SkipTLSVerify,
		},
	}
	if len(osConnection.Params.ProxyUrl) > 0 {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			parsed, err := url.Parse(osConnection.Params.ProxyUrl)
			if err != nil {
				return nil, err
			}
			return parsed, nil
		}
	}
	config := opensearch.Config{
		Transport: transport,
		Addresses: []string{osConnection.Params.Server},
		Username:  userCreds.Username,
		Password:  userCreds.Password,
	}
	if (ctx.Value(consts.DebugFlag) != nil && ctx.Value(consts.DebugFlag).(bool)) || (c.CliParams != nil && c.CliParams.EnableDebugLogs()) {
		config.EnableDebugLogger = true
	}
	return config, nil
}

// GetOpenSearchClient returns an OpenSearch client based on the given app configuration.
func GetOpenSearchClient(c appconfig.AppConfig, ctx context.Context) (*opensearch.Client, error) {
	config, osConfigErr := BuildOSConfig(c, ctx)
	if osConfigErr != nil {
		return nil, osConfigErr
	}

	if client, clientInitErr := opensearchapi.NewClient(opensearchapi.Config{Client: config}); clientInitErr != nil {
		log.Warn().Msg("unable to create client, check your config file.")
		return nil, clientInitErr
	} else {
		if clusterInfo, fetchInfoErr := client.Info(context.Background(), nil); fetchInfoErr != nil {
			log.Warn().Msg("unable to discover nodes, check your config file.")
			return nil, fetchInfoErr
		} else {
			log.Debug().Interface("clusterInfo", clusterInfo).Msg("current cluster info")
			log.Debug().Msg("client initialized")
		}
		return client.Client, nil
	}
}
