package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// GetStatsLag retrieves and displays replication lag statistics for a specified index.
// function wraps the following opensearch-go API call:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/api/#get-replication-status
func (api *OpensearchWrapper) GetStatsLag(indexName string, raw bool) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result IndexReplicationStatusResponse
	_, err := api.Client.Do(ctx, IndexReplicationStatsReq{Index: indexName, Params: IndexReplicationStatsParams{Verbose: true}}, &result)
	if err != nil {
		log.Fatal(err)
	}
	if raw {
		bytes, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("\n%s\n", bytes)
	} else {
		switch strings.ToUpper(result.Status) {
		case "SYNCING":
			fmt.Println("replication is in sync")
			fmt.Printf("lag value (follower_checkpoint - leader_checkpoint):%d", result.SyncingDetails.FollowerCheckpoint-result.SyncingDetails.LeaderCheckpoint)
		case "BOOTSTRAPING":
			fmt.Println("replication is in bootstrap mode")
		case "PAUSED":
			fmt.Println("replication is paused")
		case "REPLICATION NOT IN PROGRESS":
			fmt.Println("replication is not in progress")
		case "FAILED":
			fmt.Printf("replication failed for index '%s'\nreason:\n%s", indexName, result.Reason)

		default:
			fmt.Println("replication status is unknown")

		}
	}
}
