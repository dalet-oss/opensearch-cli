package api

import (
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"golang.org/x/exp/maps"
	"log"
	"slices"
	"strings"
)

const (
	SecurityPlugin = "opensearch-security"
	CCRPlugin      = "opensearch-cross-cluster-replication"
)

func HasPlugin(pluginsList []opensearchapi.CatPluginResp, name string) bool {
	return slices.ContainsFunc(pluginsList, func(e opensearchapi.CatPluginResp) bool {
		return strings.Contains(e.Component, name)
	})
}

func (api *OpensearchWrapper) getClusterSettings() opensearchapi.ClusterGetSettingsResp {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result opensearchapi.ClusterGetSettingsResp
	_, err := api.Client.Do(ctx, opensearchapi.ClusterGetSettingsReq{Params: opensearchapi.ClusterGetSettingsParams{Pretty: false}}, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func buildClusterRemoteDeleteMap(remoteName string, toNull map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"cluster": map[string]interface{}{
			"remote": map[string]interface{}{
				remoteName: jsonNullify(toNull),
			},
		},
	}
}

func deleteRemote(remoteName string, settings opensearchapi.ClusterGetSettingsResp) (result map[string]interface{}) {
	result = make(map[string]interface{})
	var persistentSettings map[string]interface{}
	var transientSettings map[string]interface{}
	parseErr := json.Unmarshal(settings.Persistent, &persistentSettings)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	parseErr = json.Unmarshal(settings.Transient, &transientSettings)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	transientRemote, foundInTransientSettings := findRemoteSettings(remoteName, transientSettings)
	persistentRemote, foundInPersistentSettings := findRemoteSettings(remoteName, transientSettings)

	if foundInTransientSettings || foundInPersistentSettings {
		if foundInTransientSettings {
			result["transient"] = buildClusterRemoteDeleteMap(remoteName, transientRemote)
		}
		if foundInPersistentSettings {
			result["persistent"] = buildClusterRemoteDeleteMap(remoteName, persistentRemote)
		}
	}
	return result
}

func jsonNullify(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for _, k := range maps.Keys(in) {
		out[k] = []byte("null")
	}
	return out
}

func findRemoteSettings(remoteName string, settings map[string]interface{}) (map[string]interface{}, bool) {
	if v, ok := settings["cluster"]; !ok {
		return nil, false
	} else {
		clusterSettings := v.(map[string]interface{})
		if v, hasRemote := clusterSettings["remote"]; !hasRemote {
			return nil, false
		} else {
			remoteSettings := v.(map[string]interface{})
			if remoteConfig, hasRemoteWithName := remoteSettings[remoteName]; !hasRemoteWithName {
				return nil, false
			} else {
				return remoteConfig.(map[string]interface{}), true
			}
		}
	}
}
