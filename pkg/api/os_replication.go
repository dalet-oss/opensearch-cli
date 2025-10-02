package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/replication"
	tstats "bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/stats"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"log"
)

// CreateReplication creates the replication task
// Initiate replication of an index from the leader cluster to the follower cluster. Send this request to the follower cluster.
func (api *OpensearchWrapper) CreateReplication(opts replication.StartReplicationReq, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, opts, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("create replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}

func (api *OpensearchWrapper) PauseReplication(indexName string, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, replication.PauseReplicationReq{Index: indexName}, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("pause replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}
func (api *OpensearchWrapper) ResumeReplication(indexName string, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, replication.ResumeReplicationReq{Index: indexName}, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("pause replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}
func (api *OpensearchWrapper) StopReplication(indexName string, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, replication.StopReplicationReq{Index: indexName}, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("stop replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}
func (api *OpensearchWrapper) StatusReplication(indexName string, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result tstats.IndexReplicationStatsResponse
	params := tstats.IndexReplicationStatsReq{Index: indexName, Params: tstats.IndexReplicationStatsParams{Verbose: true}}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("replication status for index '%s':\n%s\n", indexName, printutils.MarshalJSONOrDie(result))
		}
	}
}
func (api *OpensearchWrapper) TaskStatusReplication(index string, detailed, table, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	query := opensearchapi.CatRecoveryReq{Params: opensearchapi.CatRecoveryParams{
		ActiveOnly: opensearchapi.ToPointer(false),
		Detailed:   opensearchapi.ToPointer(detailed),
		V:          opensearchapi.ToPointer(true),
		Pretty:     true,
	}}
	if len(index) > 0 {
		query.Indices = []string{index}
	}
	var result []opensearchapi.CatRecoveryItemResp
	if rsp, err := api.Client.Do(ctx, query, &result); err != nil {
		log.Fatal(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Printf("recovery status:\n%s", printutils.MarshalJSONOrDie(result))
		}
	}
}
