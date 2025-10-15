package api

import (
	tstats "bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/stats"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
	"strings"
)

// GetStatsLag retrieves and displays replication lag statistics for a specified index.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-replication-status
func (api *OpensearchWrapper) GetStatsLag(indexName string, raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()
	var result tstats.IndexReplicationStatsResponse
	rsp, err := api.Client.Do(ctx, tstats.IndexReplicationStatsReq{Index: indexName, Params: tstats.IndexReplicationStatsParams{Verbose: true}}, &result)
	if err != nil {
		log.Fatal().Err(err)
	}
	if raw || rsp.IsError() {
		bytes := printutils.MarshalJSONOrDie(result)
		log.Info().Msgf("\n%s\n", bytes)
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
			log.Fatal().Msgf("replication failed for index '%s'\nreason:\n%s", indexName, result.Reason)
		default:
			log.Fatal().Msg("replication status is unknown")

		}
	}
}

// GetReplicationLeaderStats retrieves and displays replication leader statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-leader-cluster-stats
func (api *OpensearchWrapper) GetReplicationLeaderStats(raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationLeaderStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationLeaderStatsReq{}, &result)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("\n%s\n", printutils.MarshalJSONOrDie(result))
}

// GetReplicationFollowerStats retrieves and displays replication follower statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-follower-cluster-stats
func (api *OpensearchWrapper) GetReplicationFollowerStats(raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationFollowerStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationFollowerStatsReq{}, &result)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("\n%s\n", printutils.MarshalJSONOrDie(result))
}

// GetReplicationAutofollowStats retrieves and displays replication autofollow statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-auto-follow-stats
func (api *OpensearchWrapper) GetReplicationAutofollowStats(raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationAutoFollowStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationAutoFollowStatsReq{}, &result)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("\n%s\n", printutils.MarshalJSONOrDie(result))
}

// ListOfAFRules - shows the list of configured autofollow rules
func (api *OpensearchWrapper) ListOfAFRules(raw bool) {
	ctx, cancelFunc := api.requestContext()
	defer cancelFunc()

	var result tstats.ReplicationAutoFollowStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationAutoFollowStatsReq{}, &result)
	if err != nil {
		log.Fatal().Err(err)
	}
	if !raw {
		log.Info().Msg("configured autofollow rules:")
		for _, afr := range result.AutofollowStats {
			log.Info().Msgf("name: '%s' | pattern: '%s' | failed indicies: %v", afr.Name, afr.Pattern, afr.FailedIndices)
		}
	} else {
		log.Info().Msgf("configured autofollow rules: \n%s", printutils.MarshalJSONOrDie(result.AutofollowStats))
	}
}
