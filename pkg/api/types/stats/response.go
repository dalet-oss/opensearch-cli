package stats

type IndexReplicationStatsResponse struct {
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	LeaderAlias    string `json:"leader_alias"`
	LeaderIndex    string `json:"leader_index"`
	FollowerIndex  string `json:"follower_index"`
	SyncingDetails struct {
		LeaderCheckpoint   int `json:"leader_checkpoint"`
		FollowerCheckpoint int `json:"follower_checkpoint"`
		SeqNo              int `json:"seq_no"`
	} `json:"syncing_details"`
}

type ReplicationLeaderStatsResponse struct {
	NumReplicatedIndices        int                                    `json:"num_replicated_indices"`
	OperationsRead              int                                    `json:"operations_read"`
	TranslogSizeBytes           int                                    `json:"translog_size_bytes"`
	OperationsReadLucene        int                                    `json:"operations_read_lucene"`
	OperationsReadTranslog      int                                    `json:"operations_read_translog"`
	TotalReadTimeLuceneMillis   int                                    `json:"total_read_time_lucene_millis"`
	TotalReadTimeTranslogMillis int                                    `json:"total_read_time_translog_millis"`
	BytesRead                   int                                    `json:"bytes_read"`
	IndexStats                  map[string]ReplicationLeaderIndexStats `json:"index_stats"`
}

type ReplicationLeaderIndexStats struct {
	OperationsRead              int `json:"operations_read"`
	TranslogSizeBytes           int `json:"translog_size_bytes"`
	OperationsReadLucene        int `json:"operations_read_lucene"`
	OperationsReadTranslog      int `json:"operations_read_translog"`
	TotalReadTimeLuceneMillis   int `json:"total_read_time_lucene_millis"`
	TotalReadTimeTranslogMillis int `json:"total_read_time_translog_millis"`
	BytesRead                   int `json:"bytes_read"`
}

type ReplicationFollowerStatsResponse struct {
	NumSyncingIndices       int                                      `json:"num_syncing_indices"`
	NumBootstrappingIndices int                                      `json:"num_bootstrapping_indices"`
	NumPausedIndices        int                                      `json:"num_paused_indices"`
	NumFailedIndices        int                                      `json:"num_failed_indices"`
	NumShardTasks           int                                      `json:"num_shard_tasks"`
	NumIndexTasks           int                                      `json:"num_index_tasks"`
	OperationsWritten       int                                      `json:"operations_written"`
	OperationsRead          int                                      `json:"operations_read"`
	FailedReadRequests      int                                      `json:"failed_read_requests"`
	ThrottledReadRequests   int                                      `json:"throttled_read_requests"`
	FailedWriteRequests     int                                      `json:"failed_write_requests"`
	ThrottledWriteRequests  int                                      `json:"throttled_write_requests"`
	FollowerCheckpoint      int                                      `json:"follower_checkpoint"`
	LeaderCheckpoint        int                                      `json:"leader_checkpoint"`
	TotalWriteTimeMillis    int                                      `json:"total_write_time_millis"`
	IndexStats              map[string]ReplicationFollowerIndexStats `json:"index_stats"`
}

type ReplicationFollowerIndexStats struct {
	OperationsWritten      int `json:"operations_written"`
	OperationsRead         int `json:"operations_read"`
	FailedReadRequests     int `json:"failed_read_requests"`
	ThrottledReadRequests  int `json:"throttled_read_requests"`
	FailedWriteRequests    int `json:"failed_write_requests"`
	ThrottledWriteRequests int `json:"throttled_write_requests"`
	FollowerCheckpoint     int `json:"follower_checkpoint"`
	LeaderCheckpoint       int `json:"leader_checkpoint"`
	TotalWriteTimeMillis   int `json:"total_write_time_millis"`
}

type ReplicationAutofollowStatsResponse struct {
	NumSuccessStartReplication int           `json:"num_success_start_replication"`
	NumFailedStartReplication  int           `json:"num_failed_start_replication"`
	NumFailedLeaderCalls       int           `json:"num_failed_leader_calls"`
	FailedIndices              []interface{} `json:"failed_indices"`
	AutofollowStats            []struct {
		Name                       string        `json:"name"`
		Pattern                    string        `json:"pattern"`
		NumSuccessStartReplication int           `json:"num_success_start_replication"`
		NumFailedStartReplication  int           `json:"num_failed_start_replication"`
		NumFailedLeaderCalls       int           `json:"num_failed_leader_calls"`
		FailedIndices              []interface{} `json:"failed_indices"`
	} `json:"autofollow_stats"`
}
