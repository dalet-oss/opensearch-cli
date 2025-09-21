package api

import (
	"context"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"log"
)

type IndexInfo struct {
	Health       string `json:"health,omitempty"`
	Status       string `json:"status,omitempty"`
	Index        string `json:"index,omitempty"`
	Uuid         string `json:"uuid,omitempty"`
	Pri          string `json:"pri,omitempty"`
	Rep          string `json:"rep,omitempty"`
	DocsCount    string `json:"docs.count,omitempty"`
	DocsDeleted  string `json:"docs.deleted,omitempty"`
	StoreSize    string `json:"store.size,omitempty"`
	PriStoreSize string `json:"pri.store.size,omitempty"`
}

type IndexInfoResponse []IndexInfo

// GetIndexList returns a list of all indices.
// somehow the lib doesn't allow using of the _list/indices endpoint
// exposed [to the lib code] only the _cat/indices endpoint
// docs: https://docs.opensearch.org/2.19/api-reference/cat/cat-indices/
func (api *OpensearchWrapper) GetIndexList() IndexInfoResponse {
	request := opensearchapi.CatIndicesReq{Params: opensearchapi.CatIndicesParams{}}
	responseData := IndexInfoResponse{}
	_, err := api.Client.Do(context.Background(), request, &responseData)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

func (api *OpensearchWrapper) DeleteIndex(indexName string) {

}
