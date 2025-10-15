package api

import (
	"bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/replication"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
)

// CreateAutofollowRule - Automatically starts replication on indexes matching a specified pattern.
// If a new index on the leader cluster matches the pattern, OpenSearch automatically creates a follower index and begins replication.
func (api *OpensearchWrapper) CreateAutofollowRule(opts replication.CreateAutofollowReq, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result interface{}
	if rsp, err := api.Client.Do(ctx, opts, &result); err != nil {
		log.Fatal().Err(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Info().Msgf("autofollow rule creation result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}

// DeleteAutofollow - Deletes the specified replication rule.
// This operation prevents any new indexes from being replicated but does not stop existing replication that the rule has already initiated.
// Replicated indexes remain read-only until you stop replication.
func (api *OpensearchWrapper) DeleteAutofollow(opts replication.DeleteAutofollowReq, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result interface{}

	if rsp, err := api.Client.Do(ctx, opts, &result); err != nil {
		log.Fatal().Err(err)
	} else {
		if raw || rsp.IsError() {
			printutils.RawResponse(rsp)
		} else {
			log.Info().Msgf("autofollow rule deletion result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
}
