package api

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/stats"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const leaderOsInst = "leader"

// const followerOsInst = "follower"
var normalReplicationRoles = replication.ReplicationRoles{
	LeaderClusterRole:   "admin",
	FollowerClusterRole: "admin",
}

// preConfigureOSIndex ensures the provided OpenSearch instance is ready with the specified index for testing purposes.
// It validates the OpenSearch wrapper, creates the specified index, and logs the progress during configuration.
func preConfigureOSIndex(t *testing.T, c *OpensearchWrapper, osInstance, replicatedIndexName string) {
	assert.NotNil(t, c)
	t.Logf("[%s]creating index", osInstance)
	assert.NoError(t, c.CreateIndex(replicatedIndexName))
	t.Logf("[%s]configured", osInstance)
}

// getStartReplicationQuery constructs and returns a StartReplicationReq for initiating index replication on a target cluster.
// only for testing purposes
func getStartReplicationQuery(replicatedIndexName string) replication.StartReplicationReq {
	return replication.StartReplicationReq{
		Index: replicatedIndexName,
		Body: replication.StartReplicationBody{
			LeaderAlias: getCCR().RemoteName,
			LeaderIndex: replicatedIndexName,
			UseRoles:    normalReplicationRoles,
		},
	}
}

func leaderShotgunInstance(index string, delay time.Duration) *shotgun {
	return NewShotgun(wrapperForContainer(LeaderContainer), false, index, shotgunBasicDocument, true, delay)
}

// TestOpensearchWrapper_CreateReplication tests creation of the replication task on the follower cluster.
// The test is executed only if the cluster has replication plugin installed.
// replication pre-requisites:
// https://docs.opensearch.org/2.19/tuning-your-cluster/replication-plugin/getting-started/#start-replication
// normal flow outline:
//
// Leader cluster:
//   - create index
//   - populate index with data
//
// Follower:
//   - configure connection to the remote cluster
//   - create replication task
//   - confirm replication status
func TestOpensearchWrapper_CreateReplication(t *testing.T) {
	replicatedIndexName := "shotgun-create-replication-task-test-index"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|create replication task",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(100),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]configured]")
			},
			CaseInput: getStartReplicationQuery(replicatedIndexName),
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
		{
			Name:      "negative| remote cluster is not configured",
			WantErr:   true,
			Wrapper:   testWrapper(),
			CaseInput: getStartReplicationQuery(replicatedIndexName),
		},
		{
			Name:    "negative| remote cluster configured but doesn't have the index",
			WantErr: true,
			Wrapper: testWrapper(),
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]configured]")
			},
			CaseInput: replication.StartReplicationReq{
				Index: "manual-replication-start-test-index-does-not-exist",
				Body: replication.StartReplicationBody{
					LeaderAlias: getCCR().RemoteName,
					LeaderIndex: replicatedIndexName,
				},
			},
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				t.Log("[follower]cleaned up")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// pre | post setup
			t.Cleanup(func() {
				if tt.PostFollowerFunc != nil {
					tt.PostFollowerFunc(t, wrapperForContainer(MainContainer))
				}
				if tt.PostLeaderFunc != nil {
					tt.PostLeaderFunc(t, wrapperForContainer(LeaderContainer))
				}
			})
			t.Log("configuring OpenSearch instance(s)")
			if tt.ConfigureLeaderFunc != nil {
				tt.ConfigureLeaderFunc(t, wrapperForContainer(LeaderContainer))
			}
			if tt.ConfigureFollowerFunc != nil {
				tt.ConfigureFollowerFunc(t, wrapperForContainer(MainContainer))
			}
			if tt.Shotgun != nil && tt.DocumentCount != nil {
				tt.Shotgun.Shoot(t, *tt.DocumentCount, nil)
			}
			t.Log("configured")
			// actual test
			executionError := tt.Wrapper.CreateReplication(tt.CaseInput.(replication.StartReplicationReq), false)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
		})
	}
}

// TestOpensearchWrapper_PauseReplication tests pausing of the replication task on the follower cluster.
// The test is executed only if the cluster has replication plugin installed
//
// normal flow outline
//
// Leader cluster:
//   - create index
//   - populate index with some data
//
// Follower cluster:
//   - configure connection to the remote cluster
//   - create the replication task
//   - pause the replication
//   - confirm paused state
func TestOpensearchWrapper_PauseReplication(t *testing.T) {
	replicatedIndexName := "shotgun-pause-replication-task-test-index"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|pausing replication",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(50),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]creating replication task")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				if err := c.StopReplication(replicatedIndexName, true); err != nil {
					t.Log("[follower]failed to stop replication")
				}
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
		{
			Name:          "negative|pausing replication that doesn't exist",
			WantErr:       true,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(50),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// pre | post setup
			t.Log("configuring OpenSearch instance(s)")
			t.Cleanup(func() {
				if tt.PostFollowerFunc != nil {
					tt.PostFollowerFunc(t, wrapperForContainer(MainContainer))
				}
				if tt.PostLeaderFunc != nil {
					tt.PostLeaderFunc(t, wrapperForContainer(LeaderContainer))
				}
			})
			if tt.ConfigureLeaderFunc != nil {
				tt.ConfigureLeaderFunc(t, wrapperForContainer(LeaderContainer))
			}
			if tt.ConfigureFollowerFunc != nil {
				tt.ConfigureFollowerFunc(t, wrapperForContainer(MainContainer))
			}
			if tt.Shotgun != nil && tt.DocumentCount != nil {
				tt.Shotgun.Shoot(t, *tt.DocumentCount, nil)
			}
			t.Log("configured")
			// actual test
			executionError := tt.Wrapper.PauseReplication(tt.CaseInput.(string), true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
				status, replStatusErr := tt.Wrapper.StatusReplication(tt.CaseInput.(string), true)
				assert.NoError(t, replStatusErr, "expected to get no error")
				assert.Contains(t, status.Status, "PAUSED", "expected status to contain paused")
			}
		})
	}
}

// TestOpensearchWrapper_ResumeReplication tests resuming of the replication task on the follower cluster.
// The test is executed only if the cluster has replication plugin installed
//
// normal flow outline
//
// Leader cluster:
//   - create index
//   - populate index with some data
//
// Follower cluster:
//   - configure connection to the remote cluster
//   - create the replication task
//   - pause the replication
//   - confirm paused state
//   - resume replication
//   - confirm resumes state of the replication
func TestOpensearchWrapper_ResumeReplication(t *testing.T) {
	replicatedIndexName := "shotgun-resume-replication-task-test-index"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|resume replication",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(10),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]creating replication task")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				time.Sleep(1 * time.Second)
				t.Log("[follower] pausing replication")
				pauseErr := c.PauseReplication(replicatedIndexName, true)
				assert.NoError(t, pauseErr, "expected to pause replication")
				replicationTaskStatus, statusQueryErr := c.StatusReplication(replicatedIndexName, true)
				assert.NoError(t, statusQueryErr, "expected to get no error")
				assert.Equal(t, replicationTaskStatus.Status, "PAUSED", "expected status PAUSED")
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				if err := c.StopReplication(replicatedIndexName, true); err != nil {
					t.Log("[follower]failed to stop replication")
				}
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
		{
			Name:          "negative|resume replication that is not paused",
			WantErr:       true,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(10),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]creating replication task")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				if err := c.StopReplication(replicatedIndexName, true); err != nil {
					t.Log("[follower]failed to stop replication")
				}
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// pre | post setup
			t.Log("configuring OpenSearch instance(s)")
			t.Cleanup(func() {
				if tt.PostFollowerFunc != nil {
					tt.PostFollowerFunc(t, wrapperForContainer(MainContainer))
				}
				if tt.PostLeaderFunc != nil {
					tt.PostLeaderFunc(t, wrapperForContainer(LeaderContainer))
				}
			})
			if tt.ConfigureLeaderFunc != nil {
				tt.ConfigureLeaderFunc(t, wrapperForContainer(LeaderContainer))
			}
			if tt.Shotgun != nil && tt.DocumentCount != nil {
				tt.Shotgun.Shoot(t, *tt.DocumentCount, nil)
			}
			if tt.ConfigureFollowerFunc != nil {
				tt.ConfigureFollowerFunc(t, wrapperForContainer(MainContainer))
			}
			t.Log("configured")
			// actual test
			executionError := tt.Wrapper.ResumeReplication(tt.CaseInput.(string), true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
				replicationTaskStatus, statusQueryErr := tt.Wrapper.StatusReplication(tt.CaseInput.(string), true)
				assert.NoError(t, statusQueryErr, "expected to get no error")
				assert.Equal(t, replicationTaskStatus.Status, "SYNCING", "expected status to contain paused")
			}
		})
	}
}

// TestOpensearchWrapper_StatusReplication tests querying status of the replication task on the follower cluster.
// The test is executed only if the cluster has replication plugin installed
//
// normal flow outline
//
// Leader cluster:
//   - create index
//   - populate index with some data
//
// Follower cluster:
//   - configure the remote
//   - configure the replication task
//   - query status
//   - confirm the expected status of the task
func TestOpensearchWrapper_StatusReplication(t *testing.T) {
	replicatedIndexName := "shotgun-status-replication-task-test-index"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|query existing replication",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(10),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]creating replication task")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				if err := c.StopReplication(replicatedIndexName, true); err != nil {
					t.Log("[follower]failed to stop replication")
				}
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
			ExtraValidationFunc: func(t *testing.T, execResult any) {
				repl := execResult.(stats.IndexReplicationStatsResponse)
				assert.Equal(t, repl.Status, "SYNCING",
					"expected to get SYNCING")
			},
		},
		{
			// for non-existing replication OS returns: HTTP 200, "status": "REPLICATION NOT IN PROGRESS"
			Name:                  "negative|query non-existing replication ",
			WantErr:               false,
			Shotgun:               nil,
			DocumentCount:         nil,
			Wrapper:               testWrapper(),
			ConfigureLeaderFunc:   nil,
			ConfigureFollowerFunc: nil,
			CaseInput:             replicatedIndexName,
			ExtraValidationFunc: func(t *testing.T, execResult any) {
				repl := execResult.(stats.IndexReplicationStatsResponse)
				assert.Equal(t, repl.Status, "REPLICATION NOT IN PROGRESS",
					"expected to get REPLICATION NOT IN PROGRESS for non-existing replication")
			},
			PostFollowerFunc: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// pre | post setup
			t.Log("configuring OpenSearch instance(s)")
			t.Cleanup(func() {
				if tt.PostFollowerFunc != nil {
					tt.PostFollowerFunc(t, wrapperForContainer(MainContainer))
				}
				if tt.PostLeaderFunc != nil {
					tt.PostLeaderFunc(t, wrapperForContainer(LeaderContainer))
				}
			})
			if tt.ConfigureLeaderFunc != nil {
				tt.ConfigureLeaderFunc(t, wrapperForContainer(LeaderContainer))
			}
			if tt.ConfigureFollowerFunc != nil {
				tt.ConfigureFollowerFunc(t, wrapperForContainer(MainContainer))
			}
			if tt.Shotgun != nil && tt.DocumentCount != nil {
				tt.Shotgun.Shoot(t, *tt.DocumentCount, nil)
			}
			t.Log("configured")
			// actual test
			replicationStatus, executionError := tt.Wrapper.StatusReplication(tt.CaseInput.(string), true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
			if tt.ExtraValidationFunc != nil {
				tt.ExtraValidationFunc(t, replicationStatus)
			}
		})
	}
}

// TestOpensearchWrapper_StopReplication tests stopping of the replication task on the follower cluster.
// The test is executed only if the cluster has replication plugin installed
//
// normal flow outline
//
// Leader cluster:
//   - create index
//   - populate index with some data
//
// Follower cluster:
//   - configure the remote
//   - configure the replication task
//   - stop the replication task
func TestOpensearchWrapper_StopReplication(t *testing.T) {
	replicatedIndexName := "shotgun-stop-replication-task-test-index"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|stop active replication",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(10),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]creating replication task")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
		{
			Name:                  "negative|stop non-existing replication",
			WantErr:               true,
			Shotgun:               nil,
			DocumentCount:         nil,
			Wrapper:               testWrapper(),
			ConfigureLeaderFunc:   nil,
			ConfigureFollowerFunc: nil,
			CaseInput:             replicatedIndexName,
			PostFollowerFunc:      nil,
		},
		{
			Name:          "negative|stop already stopped replication",
			WantErr:       true,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(10),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("[follower]creating replication task")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				time.Sleep(1 * time.Second)
				t.Log("[follower] stopping replication")
				stopErr := c.StopReplication(replicatedIndexName, true)
				assert.NoError(t, stopErr, "expected to stop replication")
				replicationTaskStatus, statusQueryErr := c.StatusReplication(replicatedIndexName, true)
				assert.NoError(t, statusQueryErr, "expected to get no error")
				assert.Equal(t, replicationTaskStatus.Status, "REPLICATION NOT IN PROGRESS", "expected status REPLICATION NOT IN PROGRESS")
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// pre | post setup
			t.Log("configuring OpenSearch instance(s)")
			t.Cleanup(func() {
				if tt.PostFollowerFunc != nil {
					tt.PostFollowerFunc(t, wrapperForContainer(MainContainer))
				}
				if tt.PostLeaderFunc != nil {
					tt.PostLeaderFunc(t, wrapperForContainer(LeaderContainer))
				}
			})
			if tt.ConfigureLeaderFunc != nil {
				tt.ConfigureLeaderFunc(t, wrapperForContainer(LeaderContainer))
			}
			if tt.ConfigureFollowerFunc != nil {
				tt.ConfigureFollowerFunc(t, wrapperForContainer(MainContainer))
			}
			if tt.Shotgun != nil && tt.DocumentCount != nil {
				tt.Shotgun.Shoot(t, *tt.DocumentCount, nil)
			}
			t.Log("configured")
			// actual test
			executionError := tt.Wrapper.StopReplication(tt.CaseInput.(string), true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
		})
	}
}

// TestOpensearchWrapper_TaskStatusReplication performs tests for OpenSearch replication status functionality.
// It ensures proper configuration and cleanup of leader and follower indices and verifies replication task execution.
// Multiple test cases are run with pre/post setup steps, and assertions are made on the expected errors or results.
func TestOpensearchWrapper_TaskStatusReplication(t *testing.T) {
	replicatedIndexName := "shotgun-task-status-replication-task-test-index"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|query existing replication",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndexName, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(10),
			Wrapper:       testWrapper(),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				preConfigureOSIndex(t, c, leaderOsInst, replicatedIndexName)
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				assert.NoError(t, c.CreateReplication(getStartReplicationQuery(replicatedIndexName), true))
				t.Log("[follower]configured]")
			},
			CaseInput: replicatedIndexName,
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]cleaning up")
				assert.NoError(t, c.StopReplication(replicatedIndexName, true))
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				assert.NoError(t, c.DeleteRemote(getCCR().RemoteName, true))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up")
				assert.NoError(t, c.DeleteIndex(replicatedIndexName))
				t.Log("[leader]cleaned up")
			},
		},
		{
			Name:                  "negative|query non-existing replication",
			WantErr:               true,
			Shotgun:               nil,
			DocumentCount:         nil,
			Wrapper:               testWrapper(),
			ConfigureLeaderFunc:   nil,
			ConfigureFollowerFunc: nil,
			CaseInput:             replicatedIndexName,
			PostFollowerFunc:      nil,
			PostLeaderFunc:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// pre | post setup
			t.Log("configuring OpenSearch instance(s)")
			t.Cleanup(func() {
				if tt.PostFollowerFunc != nil {
					tt.PostFollowerFunc(t, wrapperForContainer(MainContainer))
				}
				if tt.PostLeaderFunc != nil {
					tt.PostLeaderFunc(t, wrapperForContainer(LeaderContainer))
				}
			})
			if tt.ConfigureLeaderFunc != nil {
				tt.ConfigureLeaderFunc(t, wrapperForContainer(LeaderContainer))
			}
			if tt.ConfigureFollowerFunc != nil {
				tt.ConfigureFollowerFunc(t, wrapperForContainer(MainContainer))
			}
			if tt.Shotgun != nil && tt.DocumentCount != nil {
				tt.Shotgun.Shoot(t, *tt.DocumentCount, nil)
			}
			t.Log("configured")
			// actual test
			executionError := tt.Wrapper.TaskStatusReplication(tt.CaseInput.(string), true, true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
		})
	}
}
