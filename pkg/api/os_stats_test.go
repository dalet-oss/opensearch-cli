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

func TestOpensearchWrapper_GetReplicationAutofollowStats(t *testing.T) {
	tests := []OSMultiContainerTest{}
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

func TestOpensearchWrapper_GetReplicationFollowerStats(t *testing.T) {
	tests := []OSMultiContainerTest{}
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

func TestOpensearchWrapper_GetReplicationLeaderStats(t *testing.T) {
	tests := []OSMultiContainerTest{}
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

func TestOpensearchWrapper_GetStatsLag(t *testing.T) {
	tests := []OSMultiContainerTest{}
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
			statsLag, executionError := tt.Wrapper.GetStatsLag(tt.CaseInput.(string), true)
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
