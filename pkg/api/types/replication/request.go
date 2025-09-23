package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v4"
	"net/http"
)

type StartReplicationReq struct {
	Header http.Header
	Index  string
	Body   StartReplicationBody
}
type StartReplicationBody struct {
	LeaderAlias string `json:"leader_alias"`
	LeaderIndex string `json:"leader_index"`
	// TODO: think of gathering this info automatically from the cluster
	// UseRoles mandatory if security plugin enabled
	UseRoles struct {
		LeaderClusterRole   string `json:"leader_cluster_role"`
		FollowerClusterRole string `json:"follower_cluster_role"`
	} `json:"use_roles,omitempty"`
}

// GetRequest returns the *http.Request that gets executed by the client
func (r StartReplicationReq) GetRequest() (*http.Request, error) {
	body, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}

	return opensearch.BuildRequest(
		"PUT",
		fmt.Sprintf("/_plugins/_replication/%s/_start", r.Index),
		bytes.NewReader(body),
		make(map[string]string),
		r.Header,
	)
}

// ---

type PauseReplicationReq struct {
	Header http.Header
	Index  string
}

// GetRequest returns the *http.Request that gets executed by the client
func (r PauseReplicationReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"POST",
		fmt.Sprintf("/_plugins/_replication/%s/_pause", r.Index),
		bytes.NewReader([]byte("{}")),
		make(map[string]string),
		r.Header,
	)
}

type StopReplicationReq struct {
	Header http.Header
	Index  string
}

// GetRequest returns the *http.Request that gets executed by the client
func (r StopReplicationReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"POST",
		fmt.Sprintf("/_plugins/_replication/%s/_stop", r.Index),
		bytes.NewReader([]byte("{}")),
		make(map[string]string),
		r.Header,
	)
}

type ResumeReplicationReq struct {
	Header http.Header
	Index  string
}

// GetRequest returns the *http.Request that gets executed by the client
func (r ResumeReplicationReq) GetRequest() (*http.Request, error) {
	return opensearch.BuildRequest(
		"POST",
		fmt.Sprintf("/_plugins/_replication/%s/_resume", r.Index),
		bytes.NewReader([]byte("{}")),
		make(map[string]string),
		r.Header,
	)
}
