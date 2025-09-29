package api

import (
	tstats "bitbucket.org/ooyalaflex/opensearch-cli/pkg/api/types/stats"
	printutils "bitbucket.org/ooyalaflex/opensearch-cli/pkg/utils/print"
	"context"
	"log"
	"strings"
)

// GetStatsLag retrieves and displays replication lag statistics for a specified index.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-replication-status
func (api *OpensearchWrapper) GetStatsLag(indexName string, raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result tstats.IndexReplicationStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationStatsReq{Index: indexName, Params: tstats.IndexReplicationStatsParams{Verbose: true}}, &result)
	if err != nil {
		log.Fatal(err)
	}
	if raw {
		bytes := printutils.MarshalJSONOrDie(result)
		log.Printf("\n%s\n", bytes)
	} else {
		log.Printf("replication status for index '%s':\n", indexName)
		switch strings.ToUpper(result.Status) {
		case "SYNCING":
			log.Println("replication is in sync")
			log.Printf("lag value (follower_checkpoint - leader_checkpoint):%d", result.SyncingDetails.FollowerCheckpoint-result.SyncingDetails.LeaderCheckpoint)
		case "BOOTSTRAPPING":
			log.Println("replication is in bootstrap mode")
			log.Printf("reason:%s", result.Reason)
			log.Printf("lag value (follower_checkpoint - leader_checkpoint):%d", result.SyncingDetails.FollowerCheckpoint-result.SyncingDetails.LeaderCheckpoint)
		case "PAUSED":
			log.Println("replication is paused")
			log.Printf("reason:%s", result.Reason)
		case "REPLICATION NOT IN PROGRESS":
			log.Println("replication is not in progress")
			log.Printf("reason:%s", result.Reason)
		case "FAILED":
			log.Fatalf("replication failed for index '%s'\nreason:\n%s", indexName, result.Reason)
		default:
			log.Fatal("replication status is unknown")

		}
	}
}

// GetReplicationLeaderStats retrieves and displays replication leader statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-leader-cluster-stats
func (api *OpensearchWrapper) GetReplicationLeaderStats(raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()

	var result tstats.ReplicationLeaderStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationLeaderStatsReq{}, &result)
	if err != nil {
		log.Fatal(err)
	}
	if raw {
		log.Printf("\n%s\n", printutils.MarshalJSONOrDie(result))
	} else {
		// TODO: implement
	}
}

// GetReplicationFollowerStats retrieves and displays replication follower statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-follower-cluster-stats
func (api *OpensearchWrapper) GetReplicationFollowerStats(raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()

	var result tstats.ReplicationFollowerStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationFollowerStatsReq{}, &result)
	if err != nil {
		log.Fatal(err)
	}
	if raw {
		log.Printf("\n%s\n", printutils.MarshalJSONOrDie(result))
	} else {
		// TODO: implement
	}
}

// GetReplicationAutofollowStats retrieves and displays replication autofollow statistics for all indices.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-auto-follow-stats
func (api *OpensearchWrapper) GetReplicationAutofollowStats(raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()

	var result tstats.ReplicationAutoFollowStatsResponse
	_, err := api.Client.Do(ctx, tstats.IndexReplicationAutoFollowStatsReq{}, &result)
	if err != nil {
		log.Fatal(err)
	}
	if raw {
		log.Printf("\n%s\n", printutils.MarshalJSONOrDie(result))
	} else {
		// TODO: implement
	}
}
