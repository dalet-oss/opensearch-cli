package api

import (
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

const (
	SecurityPlugin = "opensearch-security"
	CCRPlugin      = "opensearch-cross-cluster-replication"
)

// HasPlugin checks if a plugin with the given name exists in the provided list of plugins.
func HasPlugin(pluginsList []opensearchapi.CatPluginResp, name string) bool {
	return slices.ContainsFunc(pluginsList, func(e opensearchapi.CatPluginResp) bool {
		return strings.Contains(e.Component, name)
	})
}

// getClusterSettings retrieves the cluster settings from the OpenSearch cluster and returns the response or an error.
func (api *OpensearchWrapper) getClusterSettings() (opensearchapi.ClusterGetSettingsResp, error) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result opensearchapi.ClusterGetSettingsResp
	_, err := api.Client.Do(ctx, opensearchapi.ClusterGetSettingsReq{Params: opensearchapi.ClusterGetSettingsParams{Pretty: false}}, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// buildClusterRemoteDeleteMap constructs a map to delete the remote cluster from the OpenSearch
func buildClusterRemoteDeleteMap(remoteName string, toNull map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"cluster": map[string]interface{}{
			"remote": map[string]interface{}{
				remoteName: jsonNullify(toNull),
			},
		},
	}
}

// deleteRemote deletes the remote cluster from the OpenSearch cluster.
// It returns a map with the settings to be updated.
// If the remote cluster is not found, it returns an empty map.
// If the remote cluster is found, it returns a map with the settings to be updated.
func deleteRemote(remoteName string, settings opensearchapi.ClusterGetSettingsResp) (result map[string]interface{}) {
	result = make(map[string]interface{})
	var persistentSettings map[string]interface{}
	var transientSettings map[string]interface{}
	parseErr := json.Unmarshal(settings.Persistent, &persistentSettings)
	if parseErr != nil {
		log.Fatal().Err(parseErr)
	}
	parseErr = json.Unmarshal(settings.Transient, &transientSettings)
	if parseErr != nil {
		log.Fatal().Err(parseErr)
	}
	transientRemote, foundInTransientSettings := findRemoteSettings(remoteName, transientSettings)
	persistentRemote, foundInPersistentSettings := findRemoteSettings(remoteName, persistentSettings)

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

// jsonNullify takes a map of string keys and arbitrary values, and returns a new map with all values set to nil.
func jsonNullify(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for _, k := range maps.Keys(in) {
		out[k] = nil
	}
	return out
}

// findRemoteSettings finds the remote settings in the given settings map.
// It returns the remote settings map if found, or nil if not found.
// It returns a boolean indicating whether the remote settings were found.
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
