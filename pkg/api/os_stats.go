package api

import (
	"errors"
	"fmt"
	tstats "github.com/dalet-oss/opensearch-cli/pkg/api/types/stats"
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"strings"
)

// GetStatsLag retrieves and displays replication lag statistics for a specified index.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-replication-status
func (api *OpensearchWrapper) GetStatsLag(indexName string, raw bool) (tstats.IndexReplicationStatsResponse, error) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result tstats.IndexReplicationStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationStatsReq{Index: indexName, Params: tstats.IndexReplicationStatsParams{Verbose: true}}, &result)
	if err != nil {
		return result, err
	}
	if rsp.IsError() {
		return result, errors.New(printutils.RawResponse(rsp))
	}
	if raw {
		log.Info().Msg(printutils.RawResponse(rsp))
		return result, nil
	} else {
		log.Info().Msgf("replication status for index '%s':\n", indexName)
		switch strings.ToUpper(result.Status) {
		case "SYNCING":
			log.Info().Msg("replication is in sync")
			log.Info().Msgf("lag value (follower_checkpoint - leader_checkpoint): %d", result.SyncingDetails.FollowerCheckpoint-result.SyncingDetails.LeaderCheckpoint)
		case "BOOTSTRAPPING":
			log.Info().Msg("replication is in bootstrap mode")
			log.Info().Msgf("reason:%s", result.Reason)
			log.Info().Msgf("lag value (follower_checkpoint - leader_checkpoint): %d", result.SyncingDetails.FollowerCheckpoint-result.SyncingDetails.LeaderCheckpoint)
		case "PAUSED":
			log.Info().Msg("replication is paused")
			log.Info().Msgf("reason:%s", result.Reason)
		case "REPLICATION NOT IN PROGRESS":
			log.Info().Msg("replication is not in progress")
			log.Info().Msgf("reason:%s", result.Reason)
		case "FAILED":
			return result, errors.New(fmt.Sprintf("replication failed for index '%s'\nreason:\n%s", indexName, result.Reason))
		}
	}
	return result, nil
}

// GetReplicationLeaderStats retrieves and displays replication leader statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-leader-cluster-stats
func (api *OpensearchWrapper) GetReplicationLeaderStats(raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationLeaderStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationLeaderStatsReq{}, &result)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New(printutils.RawResponse(rsp))
	}
	if raw {
		log.Info().Msg(printutils.RawResponse(rsp))
		return nil
	} else {
		log.Info().Msgf("\n%s\n", printutils.MarshalJSONOrDie(result))
	}
	return nil
}

// GetReplicationFollowerStats retrieves and displays replication follower statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-follower-cluster-stats
func (api *OpensearchWrapper) GetReplicationFollowerStats(raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationFollowerStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationFollowerStatsReq{}, &result)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New(printutils.RawResponse(rsp))
	}
	if raw {
		log.Info().Msg(printutils.RawResponse(rsp))
		return nil
	} else {
		log.Info().Msgf("\n%s\n", printutils.MarshalJSONOrDie(result))
	}
	return nil
}

// GetReplicationAutofollowStats retrieves and displays replication autofollow statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-auto-follow-stats
func (api *OpensearchWrapper) GetReplicationAutofollowStats(raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationAutoFollowStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationAutoFollowStatsReq{}, &result)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New(printutils.RawResponse(rsp))
	}
	if raw {
		log.Info().Msg(printutils.RawResponse(rsp))
		return nil
	} else {
		log.Info().Msgf("\n%s\n", printutils.MarshalJSONOrDie(result))
	}
	return nil
}

// ListOfAFRules - shows the list of configured autofollow rules
func (api *OpensearchWrapper) ListOfAFRules(raw bool) error {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationAutoFollowStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationAutoFollowStatsReq{}, &result)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New(printutils.RawResponse(rsp))
	}
	if raw {
		log.Info().Msg(printutils.RawResponse(rsp))
		return nil
	} else {
		log.Info().Msg("configured autofollow rules:")
		for _, afr := range result.AutofollowStats {
			log.Info().Msgf("name: '%s' | pattern: '%s' | failed indices: %v", afr.Name, afr.Pattern, afr.FailedIndices)
		}
	}
	return nil
}
