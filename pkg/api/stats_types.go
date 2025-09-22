package api

import (
	"fmt"
	"github.com/opensearch-project/opensearch-go/v4"
	"net/http"
)

// IndexReplicationStatsReq represents the request for getting replication stats for a specified index.
type IndexReplicationStatsReq struct {
	Header http.Header
	Index  string
	Params IndexReplicationStatsParams
}

type IndexReplicationStatsParams struct {
	Verbose bool
}

func (p IndexReplicationStatsParams) get() map[string]string {
	params := make(map[string]string)
	if p.Verbose {
		params["verbose"] = "true"
	}
	return params
}

// GetRequest returns the *http.Request that gets executed by the client
func (r IndexReplicationStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		fmt.Sprintf("/_plugins/_replication/%s/_status", r.Index),
		nil,
		r.Params.get(),
		r.Header,
	)
}

type IndexReplicationStatsResponse struct {
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	LeaderAlias    string `json:"leader_alias"`
	LeaderIndex    string `json:"leader_index"`
	FollowerIndex  string `json:"follower_index"`
	SyncingDetails struct {
		LeaderCheckpoint   int `json:"leader_checkpoint"`
		FollowerCheckpoint int `json:"follower_checkpoint"`
		SeqNo              int `json:"seq_no"`
	} `json:"syncing_details"`
}
