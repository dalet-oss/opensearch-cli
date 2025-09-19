package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/appconfig"
	"crypto/tls"
	"fmt"
	"github.com/opensearch-project/opensearch-go"
	"net/http"
)

// BuildOSConfig constructs and returns an OpenSearch configuration based on the given app configuration.
// It retrieves the active context, cluster, and user details, and incorporates user credentials into the configuration.
// The function panics if required elements like contexts, clusters, or users are missing or invalid.
func BuildOSConfig(c appconfig.AppConfig) opensearch.Config {
	var ccfg *appconfig.ContextConfig
	ccfg = c.GetActiveContext()
	if ccfg == nil {
		panic(fmt.Sprintf("context config is not found, active context is set to '%s'", c.Current))
	}
	osConnection := ccfg.GetCluster(c)
	if osConnection == nil {
		panic(fmt.Sprintf("cluster definition '%s' is not found", ccfg.Cluster))
	}
	osUser := ccfg.GetUser(c)
	if osUser == nil {
		panic(fmt.Sprintf("user definition '%s' is not found", ccfg.User))
	}
	userCreds, err := osUser.GetUserCredentials()
	if err != nil {
		panic(err)
	}
	config := opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !osConnection.Params.Tls,
			},
		},
		Addresses: []string{osConnection.Params.Server},
		Username:  userCreds.Username,
		Password:  userCreds.Password,
	}
	return config
}

// GetOpenSearchClient returns an OpenSearch client based on the given app configuration.
func GetOpenSearchClient(c appconfig.AppConfig) *opensearch.Client {
	config := BuildOSConfig(c)
	client, err := opensearch.NewClient(config)
	if err != nil {
		// todo: handle error
		panic(err)
	}
	return client
}
