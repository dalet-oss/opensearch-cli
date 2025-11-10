package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
