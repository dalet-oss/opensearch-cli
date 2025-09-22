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

//

type IndexReplicationLeaderStatsReq struct {
	Header http.Header
}

func (r IndexReplicationLeaderStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		"/_plugins/_replication/leader_stats",
		nil,
		make(map[string]string),
		r.Header,
	)
}

type IndexReplicationFollowerStatsReq struct {
	Header http.Header
}

func (r IndexReplicationFollowerStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		"/_plugins/_replication/follower_stats",
		nil,
		make(map[string]string),
		r.Header,
	)
}

type IndexReplicationAutoFollowStatsReq struct {
	Header http.Header
}

func (r IndexReplicationAutoFollowStatsReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"GET",
		"/_plugins/_replication/autofollow_stats",
		nil,
		make(map[string]string),
		r.Header,
	)
}
