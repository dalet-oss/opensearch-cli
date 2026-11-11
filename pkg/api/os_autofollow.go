package api

import (
	"errors"
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	tstats "github.com/dalet-oss/opensearch-cli/pkg/api/types/stats"
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"slices"
)

// CreateAutofollowRule - Automatically starts replication on indexes matching a specified pattern.
// If a new index on the leader cluster matches the pattern, OpenSearch automatically creates a follower index and begins replication.
func (api *OpensearchWrapper) CreateAutofollowRule(opts replication.CreateAutofollowReq, raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result any
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
			log.Info().Msgf("autofollow rule creation result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

// DeleteAutofollow - Deletes the specified replication rule.
// This operation prevents any new indexes from being replicated but does not stop existing replication that the rule has already initiated.
// Replicated indexes remain read-only until you stop replication.
func (api *OpensearchWrapper) DeleteAutofollow(opts replication.DeleteAutofollowReq, raw bool) error {
	// OpenSearch 3.x acknowledges deletion of non-existent rules and missing leader alias, so validate on our side
	if opts.Body.LeaderAlias == "" {
		return errors.New("leader alias is required to delete autofollow rule")
	}
	if exists, err := api.autofollowRuleExists(opts.Body.Name); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("no autofollow rule found with name '%s'", opts.Body.Name)
	}

	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result any

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
			log.Info().Msgf("autofollow rule deletion result:\n%s\n", printutils.MarshalJSONOrDie(result))
		}
	}
	return nil
}

// autofollowRuleExists checks whether an autofollow rule with the given name is configured.
func (api *OpensearchWrapper) autofollowRuleExists(name string) (bool, error) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationAutoFollowStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationAutoFollowStatsReq{}, &result)
	if err != nil {
		return false, err
	}
	if rsp.IsError() {
		return false, errors.New(printutils.RawResponse(rsp))
	}
	return slices.ContainsFunc(result.AutofollowStats, func(afr tstats.AutoFollowStats) bool {
		return afr.Name == name
	}), nil
}
