package api

import (
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/stats"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	gu "github.com/dalet-oss/opensearch-cli/pkg/utils/generic"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestOpensearchWrapper_GetReplicationAutofollowStats testing autofollow stats gathering from the OS instance
func TestOpensearchWrapper_GetReplicationAutofollowStats(t *testing.T) {
	replicatedIndex := "tc-replication-autofollow-stats-test"
	ccrName := "autofollow-stats-test-remote"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|get autofollow stats for perfectly configured AF",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndex, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(100),
			Wrapper:       wrapperForContainer(MainContainer),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]creating index")
				assert.NoError(t, c.CreateIndex(replicatedIndex))
				t.Log("[leader]configured]")
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getNamedCCR(ccrName), true), "expected to configure remote cluster")
				t.Log("[follower]creating af rules")
				afRule := replication.CreateAutofollowReq{
					Header: nil,
					Body: replication.CreateAutofollowBody{
						Name:         replicatedIndex,
						LeaderAlias:  ccrName,
						IndexPattern: replicatedIndex,
					},
				}
				assert.NoError(t, c.CreateAutofollowRule(afRule, true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("[follower]cleaning up")
				assert.NotNil(t, c)
				assert.NoError(t, c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name:        replicatedIndex,
						LeaderAlias: ccrName,
					},
				}, true))
				t.Log("[follower]cleaning up remote cluster")
				assert.NoError(t, c.DeleteRemote(ccrName, true))
				t.Log("[follower]stop replication")
				if err := c.StopReplication(replicatedIndex, true); err != nil {
					t.Log(err)
				}
				t.Log("[follower]cleaning up indices")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up indexes")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[leader]cleaned up")
			},
			ExtraValidationFunc: func(t *testing.T, execResult any) {
				afStats := execResult.(stats.ReplicationAutoFollowStatsResponse)
				afData := fp.Filter(afStats.AutofollowStats, func(afRuleStatistics stats.AutoFollowStats) bool {
					return afRuleStatistics.Name == replicatedIndex && afRuleStatistics.Pattern == replicatedIndex
				})
				assert.NotEmpty(t, afStats.AutofollowStats, "af stats expected to be generated")
				assert.NotEmpty(t, afData, "af rule expected to be presented")
				assert.Equal(t, 1, afData[0].NumSuccessStartReplication, "af stats expected to be generated")
			},
		},
		{
			Name:    "negative|get autofollow stats for AF with no remote cluster",
			WantErr: false,
			Wrapper: wrapperForContainer(LeaderContainer),
			ExtraValidationFunc: func(t *testing.T, execResult any) {
				afStats := execResult.(stats.ReplicationAutoFollowStatsResponse)
				assert.Empty(t, afStats.AutofollowStats, "af stats expected to be empty")
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
			autofollowStats, executionError := tt.Wrapper.GetReplicationAutofollowStats(true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
			if tt.ExtraValidationFunc != nil {
				tt.ExtraValidationFunc(t, autofollowStats)
			}
		})
	}
}

// TestOpensearchWrapper_GetReplicationFollowerStats tests the retrieval of replication follower statistics in a multi-container setup.
func TestOpensearchWrapper_GetReplicationFollowerStats(t *testing.T) {
	replicatedIndex := "tc-stats-follower-replicated-index"
	ccrName := "tc-stats-repl-follower-ccr"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|get follower stats for perfectly configured AF",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndex, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(100),
			Wrapper:       wrapperForContainer(MainContainer),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]creating index")
				assert.NoError(t, c.CreateIndex(replicatedIndex))
				t.Log("[leader]configured]")
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getNamedCCR(ccrName), true), "expected to configure remote cluster")
				t.Log("[follower]creating af rules")
				afRule := replication.CreateAutofollowReq{
					Header: nil,
					Body: replication.CreateAutofollowBody{
						Name:         replicatedIndex,
						LeaderAlias:  ccrName,
						IndexPattern: replicatedIndex,
					},
				}
				assert.NoError(t, c.CreateAutofollowRule(afRule, true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("[follower]cleaning up")
				assert.NotNil(t, c)
				assert.NoError(t, c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name:        replicatedIndex,
						LeaderAlias: ccrName,
					},
				}, true))
				t.Log("[follower]cleaning up remote cluster")
				assert.NoError(t, c.DeleteRemote(ccrName, true))
				t.Log("[follower]stop replication")
				if err := c.StopReplication(replicatedIndex, true); err != nil {
					t.Log(err)
				}
				t.Log("[follower]cleaning up indices")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up indexes")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[leader]cleaned up")
			},
			ExtraValidationFunc: nil,
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
			time.Sleep(2 * time.Second)
			// actual test
			followerStats, executionError := tt.Wrapper.GetReplicationFollowerStats(true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
			if tt.ExtraValidationFunc != nil {
				tt.ExtraValidationFunc(t, followerStats)
			}
		})
	}
}

// TestOpensearchWrapper_GetReplicationLeaderStats tests the retrieval of replication leader statistics under various scenarios.
func TestOpensearchWrapper_GetReplicationLeaderStats(t *testing.T) {
	replicatedIndex := "tc-stats-leader-replicated-index"
	ccrName := "tc-stats-repl-leader-ccr"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|get leader stats for perfectly configured AF",
			WantErr:       false,
			Shotgun:       leaderShotgunInstance(replicatedIndex, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(100),
			Wrapper:       wrapperForContainer(LeaderContainer),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]creating index")
				assert.NoError(t, c.CreateIndex(replicatedIndex))
				t.Log("[leader]configured]")
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getNamedCCR(ccrName), true), "expected to configure remote cluster")
				t.Log("[follower]creating af rules")
				afRule := replication.CreateAutofollowReq{
					Header: nil,
					Body: replication.CreateAutofollowBody{
						Name:         replicatedIndex,
						LeaderAlias:  ccrName,
						IndexPattern: replicatedIndex,
					},
				}
				assert.NoError(t, c.CreateAutofollowRule(afRule, true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("[follower]cleaning up")
				assert.NotNil(t, c)
				assert.NoError(t, c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name:        replicatedIndex,
						LeaderAlias: ccrName,
					},
				}, true))
				t.Log("[follower]cleaning up remote cluster")
				assert.NoError(t, c.DeleteRemote(ccrName, true))
				t.Log("[follower]stop replication")
				if err := c.StopReplication(replicatedIndex, true); err != nil {
					t.Log(err)
				}
				t.Log("[follower]cleaning up indices")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up indexes")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[leader]cleaned up")
			},
			ExtraValidationFunc: nil,
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
			time.Sleep(2 * time.Second)
			// actual test
			leaderStats, executionError := tt.Wrapper.GetReplicationLeaderStats(true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
			if tt.ExtraValidationFunc != nil {
				tt.ExtraValidationFunc(t, leaderStats)
			}
		})
	}
}

// TestOpensearchWrapper_GetStatsLag tests the GetStatsLag function for correct behavior in various multi-container scenarios.
func TestOpensearchWrapper_GetStatsLag(t *testing.T) {
	replicatedIndex := "tc-stats-lag-replicated-index"
	ccrName := "tc-stats-lag-ccr"
	tests := []OSMultiContainerTest{
		{
			Name:          "positive|get leader stats for perfectly configured AF",
			WantErr:       false,
			CaseInput:     replicatedIndex,
			Shotgun:       leaderShotgunInstance(replicatedIndex, 10*time.Millisecond),
			DocumentCount: fp.AsPointer(100),
			Wrapper:       wrapperForContainer(MainContainer),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]creating index")
				assert.NoError(t, c.CreateIndex(replicatedIndex))
				t.Log("[leader]configured]")
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getNamedCCR(ccrName), true), "expected to configure remote cluster")
				t.Log("[follower]creating af rules")
				afRule := replication.CreateAutofollowReq{
					Header: nil,
					Body: replication.CreateAutofollowBody{
						Name:         replicatedIndex,
						LeaderAlias:  ccrName,
						IndexPattern: replicatedIndex,
					},
				}
				assert.NoError(t, c.CreateAutofollowRule(afRule, true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("[follower]cleaning up")
				assert.NotNil(t, c)
				assert.NoError(t, c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name:        replicatedIndex,
						LeaderAlias: ccrName,
					},
				}, true))
				t.Log("[follower]cleaning up remote cluster")
				assert.NoError(t, c.DeleteRemote(ccrName, true))
				t.Log("[follower]stop replication")
				if err := c.StopReplication(replicatedIndex, true); err != nil {
					t.Log(err)
				}
				t.Log("[follower]cleaning up indices")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up indexes")
				assert.NoError(t, c.DeleteIndex(replicatedIndex))
				t.Log("[leader]cleaned up")
			},
			ExtraValidationFunc: nil,
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
			time.Sleep(2 * time.Second)
			// actual test
			statsLag, executionError := tt.Wrapper.GetStatsLag(tt.CaseInput.(string), false)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
			if tt.ExtraValidationFunc != nil {
				tt.ExtraValidationFunc(t, statsLag)
			}
		})
	}
}

// TestOpensearchWrapper_ListOfAFRules is a test function to validate the behavior of ListOfAFRules in OpensearchWrapper.
// It sets up and configures leader and follower containers to simulate auto-follow rule generation and replication.
// It ensures the expected auto-follow rules are created and validates cleanup operations post-test execution.
func TestOpensearchWrapper_ListOfAFRules(t *testing.T) {
	afRuleName := "tc-stats-af-rule"
	const afIndexPattern = "tc-stats-af-*"
	const ccrName = "tc-stats-af"
	leaderIndexCount := 10
	tests := []OSMultiContainerTest{
		{
			// this case is checked against the leader because we're never configuring AF there
			Name:    "no af rules generated",
			WantErr: false,
			Wrapper: wrapperForContainer(LeaderContainer),
			ExtraValidationFunc: func(t *testing.T, execResult any) {
				afRulesResp := execResult.(stats.ReplicationAutoFollowStatsResponse)
				assert.Empty(t, afRulesResp.AutofollowStats, "af rules expected to be emtpy")
			},
		},
		{
			Name:    "af rules generated",
			WantErr: false,
			Wrapper: wrapperForContainer(MainContainer),
			ConfigureLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]creating index")
				for i := 0; i < leaderIndexCount; i++ {
					indexName := fmt.Sprintf("%s-%d", afIndexPattern[:(len(afIndexPattern)-2)], i)
					t.Logf("[leader]creating index %s", indexName)
					assert.NoError(t, c.CreateIndex(indexName), "")
				}
				//assert.NoError(t, c.CreateIndex(afIndexPattern))
				t.Log("[leader]configured]")
			},
			ConfigureFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[follower]configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getNamedCCR(ccrName), true), "expected to configure remote cluster")
				t.Log("[follower]creating af rules")
				assert.NoError(t, c.CreateAutofollowRule(replication.CreateAutofollowReq{
					Header: nil,
					Body: replication.CreateAutofollowBody{
						Name:         afRuleName,
						LeaderAlias:  ccrName,
						IndexPattern: afIndexPattern,
					},
				}, true))
				time.Sleep(1 * time.Second)
				t.Log("[follower]configured]")
			},
			ExtraValidationFunc: func(t *testing.T, execResult any) {
				afRulesResp := execResult.(stats.ReplicationAutoFollowStatsResponse)
				assert.NotEmpty(t, afRulesResp.AutofollowStats, "af rules expected to be generated")
				assert.NotEmpty(t,
					fp.Filter(afRulesResp.AutofollowStats, func(afRuleStatistics stats.AutoFollowStats) bool {
						return afRuleStatistics.Name == afRuleName && afRuleStatistics.Pattern == afIndexPattern
					}), "af rule expected to be presented")
			},
			PostFollowerFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("[follower]cleaning up")
				assert.NotNil(t, c)
				assert.NoError(t, c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name:        afRuleName,
						LeaderAlias: ccrName,
					},
				}, true))
				t.Log("[follower]cleaning up remote cluster")
				assert.NoError(t, c.DeleteRemote(ccrName, true))
				t.Log("[follower]cleaning up indices")
				registeredIndices, err := c.GetIndexList()
				assert.NoError(t, err)
				indexNames := fp.Map(registeredIndices, func(info IndexInfo) string {
					return info.Index
				})
				for _, index := range fp.Filter(indexNames, gu.GetMatchFunc(afIndexPattern)) {
					assert.NoError(t, c.StopReplication(index, true))
					time.Sleep(100 * time.Millisecond)
					assert.NoError(t, c.DeleteIndex(index))
				}

				t.Log("[follower]cleaned up")
			},
			PostLeaderFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("[leader]cleaning up indexes")
				registeredIndices, err := c.GetIndexList()
				assert.NoError(t, err)
				indexNames := fp.Map(registeredIndices, func(info IndexInfo) string {
					return info.Index
				})
				for _, index := range fp.Filter(indexNames, gu.GetMatchFunc(afIndexPattern)) {
					assert.NoError(t, c.DeleteIndex(index))
				}
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
			afRules, executionError := tt.Wrapper.ListOfAFRules(true)
			if tt.WantErr {
				assert.Error(t, executionError, "expected to get error")
			} else {
				assert.NoError(t, executionError, "expected to get no error")
			}
			if tt.ExtraValidationFunc != nil {
				tt.ExtraValidationFunc(t, afRules)
			}
		})
	}
}
