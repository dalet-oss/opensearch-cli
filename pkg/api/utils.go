package api

import (
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
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
