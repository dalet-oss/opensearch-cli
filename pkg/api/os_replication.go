package api

import (
	"errors"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	tstats "github.com/dalet-oss/opensearch-cli/pkg/api/types/stats"
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// CreateReplication creates the replication task
// Initiate replication of an index from the leader cluster to the follower cluster. Send this request to the follower cluster.
func (api *OpensearchWrapper) CreateReplication(opts replication.StartReplicationReq, raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, opts, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("create replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

// PauseReplication pauses the replication for the specified index in OpenSearch.
// indexName specifies the name of the index whose replication is to be paused.
// raw determines whether the raw response from the API call should be logged.
// Returns an error if the API request fails or the response indicates an error.
func (api *OpensearchWrapper) PauseReplication(indexName string, raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, replication.PauseReplicationReq{Index: indexName}, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("pause replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

// ResumeReplication resumes the replication for the specified index in OpenSearch.
// indexName specifies the name of the index whose replication is to be resumed.
// raw determines whether the raw response from the API call should be logged.
// Returns an error if the API request fails or the response indicates an error.
func (api *OpensearchWrapper) ResumeReplication(indexName string, raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, replication.ResumeReplicationReq{Index: indexName}, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("pause replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

// StopReplication stops the replication process for the specified index in OpenSearch.
// It optionally logs the raw response if the raw parameter is true. Returns an error if the operation fails.
func (api *OpensearchWrapper) StopReplication(indexName string, raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}
	if rsp, err := api.Client.Do(ctx, replication.StopReplicationReq{Index: indexName}, &result); err != nil {
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("stop replication result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

// StatusReplication retrieves the replication status of a specified index in OpenSearch.
// It fetches detailed replication statistics and logs the response based on verbosity.
// Returns the IndexReplicationStatsResponse and an error, if any occurred during the operation.
func (api *OpensearchWrapper) StatusReplication(indexName string, raw bool) (tstats.IndexReplicationStatsResponse, error) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result tstats.IndexReplicationStatsResponse
	params := tstats.IndexReplicationStatsReq{Index: indexName, Params: tstats.IndexReplicationStatsParams{Verbose: true}}
	if rsp, err := api.Client.Do(ctx, params, &result); err != nil {
		return result, err
	} else {
		if rsp.IsError() {
			return result, errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return result, nil
		} else {
			log.Info().Msgf("replication status for index '%s':\n%s\n", indexName, printutils.MarshalJSONOrDie(result))
		}
	}
	return result, nil
}

// TaskStatusReplication retrieves the recovery status of a specified index in OpenSearch with optional detailed and raw outputs.
// It uses the CatRecovery API to fetch recovery details and logs the results based on the raw parameter.
func (api *OpensearchWrapper) TaskStatusReplication(index string, detailed, raw bool) error {
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
		return err
	} else {
		if rsp.IsError() {
			return errors.New(printutils.RawResponse(rsp))
		}
		if raw {
			log.Info().Msg(printutils.RawResponse(rsp))
			return nil
		} else {
			log.Info().Msgf("recovery status:\n%s", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}
