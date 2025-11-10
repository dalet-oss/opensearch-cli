package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/opensearch-project/opensearch-go/v4"
	"net/http"
)

// BuildOSConfig constructs and returns an OpenSearch configuration based on the given app configuration.
// It retrieves the active context, cluster, and user details, and incorporates user credentials into the configuration.
// The function panics if required elements like contexts, clusters, or users are missing or invalid.
func BuildOSConfig(c appconfig.AppConfig, ctx context.Context) (opensearch.Config, error) {
	var ccfg *appconfig.ContextConfig
	ccfg = c.GetActiveContext()
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
	config := opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: osConnection.Params.Tls,
			},
		},
		Addresses: []string{osConnection.Params.Server},
		Username:  userCreds.Username,
		Password:  userCreds.Password,
	}
	return config, nil
}

// GetOpenSearchClient returns an OpenSearch client based on the given app configuration.
func GetOpenSearchClient(c appconfig.AppConfig, ctx context.Context) (*opensearch.Client, error) {
	config, osConfigErr := BuildOSConfig(c, ctx)
	if osConfigErr != nil {
		return nil, osConfigErr
	}
	if client, err := opensearch.NewClient(config); err != nil {
		log.Warn().Msg("unable to create client, check your config file.")
		return nil, err
	} else {
		if discoveyErr := client.DiscoverNodes(); discoveyErr != nil {
			log.Warn().Msg("unable to discover nodes, check your config file.")
			return nil, discoveyErr
		}
		return client, nil
	}
}
