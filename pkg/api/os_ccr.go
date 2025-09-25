package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/fp"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"log"
	"strings"
)

// CCRCreateOpts - follows guidelines described at:
// https://docs.opensearch.org/2.19/install-and-configure/configuring-opensearch/index/#updating-cluster-settings-using-the-api
type CCRCreateOpts struct {
	// Type represents settings type, allowed values transient, persistent, default.
	//If empty string supplied it will be converted to persistent
	Type string
	// Mode connection mode. Available values proxy,sniff https://docs.opensearch.org/latest/install-and-configure/configuring-opensearch/cluster-settings/#cluster-level-remote-cluster-settings
	Mode string
	// RemoteName alias of the remote cluster
	RemoteName string
	// RemoteAddr address of the remote cluster.
	//Specifies the proxy server address for connecting to the remote cluster. All remote connections are routed through this single proxy endpoint.
	RemoteAddr string
}

func (opts CCRCreateOpts) BuildCCRParams() []byte {
	settings := map[string]interface{}{
		fp.GetOrDefault(opts.Type, "persistent", fp.NotEmptyString): map[string]interface{}{
			"cluster": map[string]interface{}{
				"remote": map[string]interface{}{
					opts.RemoteName: map[string]interface{}{
						"mode":          fp.GetOrDefault(opts.Mode, "proxy", fp.NotEmptyString),
						"proxy_address": opts.RemoteAddr,
					},
				},
			},
		},
	}
	return printutils.MarshalJSONOrDie(settings)
}

func (api *OpensearchWrapper) ConfigureRemoteCluster(opts CCRCreateOpts, raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result opensearchapi.ClusterPutSettingsResp
	params := opensearchapi.ClusterPutSettingsReq{
		Body:   strings.NewReader(string(opts.BuildCCRParams())),
		Params: opensearchapi.ClusterPutSettingsParams{Pretty: false},
	}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("Cross-cluster replication creation result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}

func (api *OpensearchWrapper) GetRemoteSettings(raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result opensearchapi.ClusterGetSettingsResp
	params := opensearchapi.ClusterGetSettingsReq{Params: opensearchapi.ClusterGetSettingsParams{Pretty: false}}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw {
			printutils.RawResponse(rsp)
		} else {
			var persistentSettings map[string]interface{}
			_ = json.Unmarshal(result.Persistent, &persistentSettings)
			log.Printf("Cluster remote settings:\n%s\n", printutils.MarshalJSONOrDie(persistentSettings["cluster"].(map[string]interface{})["remote"]))
		}
	}
}

func (api *OpensearchWrapper) DeleteRemote(remoteName string, raw bool) {

	remoteDeleteSettings := deleteRemote(remoteName, api.getClusterSettings())
	if len(remoteDeleteSettings) == 0 {
		log.Fatalf("No settings found for remote with name %s in the cluster", remoteName)
	}

	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result opensearchapi.ClusterPutSettingsResp
	params := opensearchapi.ClusterPutSettingsReq{
		Body:   strings.NewReader(string(printutils.MarshalJSONOrDie(remoteDeleteSettings))),
		Params: opensearchapi.ClusterPutSettingsParams{Pretty: false},
	}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("Delete remote result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}
