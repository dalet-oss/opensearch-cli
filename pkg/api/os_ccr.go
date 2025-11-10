package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"strings"
)

// CCRCreateOpts - follows guidelines described at:
// https://docs.opensearch.org/2.19/install-and-configure/configuring-opensearch/index/#updating-cluster-settings-using-the-api
type CCRCreateOpts struct {
	// Type represents settings type, allowed values transient, persistent, default.
	//If an empty string is supplied, it will be converted to persistent
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

func (api *OpensearchWrapper) ConfigureRemoteCluster(opts CCRCreateOpts, raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result opensearchapi.ClusterPutSettingsResp
	params := opensearchapi.ClusterPutSettingsReq{
		Body:   strings.NewReader(string(opts.BuildCCRParams())),
		Params: opensearchapi.ClusterPutSettingsParams{Pretty: false},
	}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("Cross-cluster replication creation result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

func (api *OpensearchWrapper) GetRemoteSettings(raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result opensearchapi.ClusterGetSettingsResp
	params := opensearchapi.ClusterGetSettingsReq{Params: opensearchapi.ClusterGetSettingsParams{Pretty: false}}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			var persistentSettings map[string]interface{}
			var clusterRemoteSettings interface{}
			_ = json.Unmarshal(result.Persistent, &persistentSettings)
			if clusterSettings, ok := persistentSettings["cluster"]; ok {
				if clusterSettingsMap, canCast := clusterSettings.(map[string]interface{}); canCast {
					if rSettings, found := clusterSettingsMap["remote"]; found {
						clusterRemoteSettings = rSettings
					} else {
						return fmt.Errorf("remote settings is not found in the cluster settings,settings:\n%v", printutils.MarshalJSONOrDie(clusterSettings))
					}
				} else {
					return fmt.Errorf("cluster settings is not map, can't cast to map[string]interface{}:\n%s", printutils.MarshalJSONOrDie(clusterSettings))
				}

			} else {
				return fmt.Errorf("cluster configuration is not found in the persistent settings,settings:\n%v", result.Persistent)
			}
			log.Info().Msgf("Cluster remote settings:\n%s\n", printutils.MarshalJSONOrDie(clusterRemoteSettings))
		}
	}
	return nil
}

func (api *OpensearchWrapper) DeleteRemote(remoteName string, raw bool) error {
	settings, err := api.getClusterSettings()
	if err != nil {
		return err
	}
	remoteDeleteSettings := deleteRemote(remoteName, settings)
	if len(remoteDeleteSettings) == 0 {
		return fmt.Errorf("no remote found with name '%s'", remoteName)
	}

	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result opensearchapi.ClusterPutSettingsResp
	params := opensearchapi.ClusterPutSettingsReq{
		Body:   strings.NewReader(string(printutils.MarshalJSONOrDie(remoteDeleteSettings))),
		Params: opensearchapi.ClusterPutSettingsParams{Pretty: false},
	}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New("delete remote error:" + printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("Delete remote result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}
