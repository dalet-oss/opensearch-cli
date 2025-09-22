package stats

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

// IndexReplicationStatsParams represents the parameters for getting replication stats for a specified index.
type IndexReplicationStatsParams struct {
	Verbose bool
}

// get returns the map of query parameters for the request.
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

// IndexReplicationLeaderStatsReq is a request type for fetching replication leader statistics for all indices.
type IndexReplicationLeaderStatsReq struct {
	Header http.Header
}

// GetRequest returns the *http.Request that gets executed by the client
func (r IndexReplicationLeaderStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		"/_plugins/_replication/leader_stats",
		nil,
		make(map[string]string),
		r.Header,
	)
}

// IndexReplicationFollowerStatsReq is a request type for fetching replication follower statistics for all indices.
type IndexReplicationFollowerStatsReq struct {
	Header http.Header
}

// GetRequest returns the *http.Request that gets executed by the client
func (r IndexReplicationFollowerStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		"/_plugins/_replication/follower_stats",
		nil,
		make(map[string]string),
		r.Header,
	)
}

// IndexReplicationAutoFollowStatsReq is a request type for fetching replication auto-follow statistics for all indices.
type IndexReplicationAutoFollowStatsReq struct {
	Header http.Header
}

// GetRequest returns the *http.Request that gets executed by the client
func (r IndexReplicationAutoFollowStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		"/_plugins/_replication/autofollow_stats",
		nil,
		make(map[string]string),
		r.Header,
	)
}
