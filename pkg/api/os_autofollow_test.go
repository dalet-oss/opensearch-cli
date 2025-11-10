package api

import (
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpensearchWrapper_CreateAutofollowRule(t *testing.T) {
	autoFollowRequest := replication.CreateAutofollowReq{
		Header: nil,
		Body: replication.CreateAutofollowBody{
			Name:         "test-create-autofollow",
			LeaderAlias:  leaderClusterName,
			IndexPattern: "test-index",
			UseRoles: replication.ReplicationRoles{
				LeaderClusterRole:   "admin",
				FollowerClusterRole: "admin",
			},
		},
	}
	// we need to skip the test if the CCR plugin is not installed, because we can't create autofollow rules without it
	for _, containerName := range osContainers {
		w := wrapperForContainer(containerName)
		plugins, err := w.PluginsList()
		assert.NoError(t, err, "expected to get plugins list")
		if !HasPlugin(plugins, CCRPlugin) {
			t.Skip("CCR plugin is not installed | .......... SKIPPED ..........")
		}
	}
	tests := []TestCase{
		{
			Name:          "create autofollow on the fresh server(no ccr configured)",
			Wrapper:       testWrapper(),
			CaseInput:     autoFollowRequest,
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:      "create autofollow rule",
			Wrapper:   testWrapper(),
			CaseInput: testIndexAutofollowRule,
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				// we need to configure the remote cluster before we can create the autofollow rule
				t.Log("configuring remote cluster")
				ccr := getCCR()
				t.Logf("remote cluster settings:%v", ccr)
				assert.NotNil(t, c)
				assert.NoError(t, c.ConfigureRemoteCluster(ccr, true), "expected to configure remote cluster")
			},
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("cleaning func")
				if err := c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name: autoFollowRequest.Body.Name,
					},
				}, true); err != nil {
					t.Logf("failed to delete autofollow rule:%v", err)
				}
				if err := c.DeleteRemote(getCCR().RemoteName, true); err != nil {
					t.Logf("failed to delete remote cluster:%v", err)
				}
			},
			WantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Cleanup(func() {
				if tt.PostFunc != nil {
					tt.PostFunc(t, tt.Wrapper)
				}
			})
			if tt.ConfigureFunc != nil {
				log.Info().Msg("executing configure func")
				tt.ConfigureFunc(t, tt.Wrapper)
			}

			executionErr := tt.Wrapper.CreateAutofollowRule(tt.CaseInput.(replication.CreateAutofollowReq), true)
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}

func TestOpensearchWrapper_DeleteAutofollow(t *testing.T) {
	autoFollowRequest := replication.CreateAutofollowReq{
		Header: nil,
		Body: replication.CreateAutofollowBody{
			Name:         "test-delete-autofollow",
			LeaderAlias:  leaderClusterName,
			IndexPattern: "test-index",
			UseRoles: replication.ReplicationRoles{
				LeaderClusterRole:   "admin",
				FollowerClusterRole: "admin",
			},
		},
	}
	tests := []TestCase{
		{
			Name:    "delete autofollow rule that doesn't exist",
			Wrapper: testWrapper(),
			CaseInput: replication.DeleteAutofollowReq{
				Header: nil,
				Body: replication.DeleteAutofollowBody{
					Name: "rule-doesn't-exist",
				},
			},
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:    "delete autofollow rule | but don't set leader alias",
			Wrapper: testWrapper(),
			CaseInput: replication.DeleteAutofollowReq{
				Header: nil,
				Body: replication.DeleteAutofollowBody{
					Name: autoFollowRequest.Body.Name,
				},
			},
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("creating autofollow rule")
				assert.NoError(t, c.CreateAutofollowRule(autoFollowRequest, true), "expected to create autofollow rule")
				t.Log("pre-configuration completed")
			},
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("cleaning for case: delete autofollow rule | but don't set leader alias")
				ccr := getCCR()
				t.Log("deleting autofollow rule")
				if err := c.DeleteAutofollow(replication.DeleteAutofollowReq{
					Header: nil,
					Body: replication.DeleteAutofollowBody{
						Name:        autoFollowRequest.Body.Name,
						LeaderAlias: ccr.RemoteName,
					},
				}, true); err != nil {
					t.Logf("failed to delete autofollow rule:%v", err)
				}
				t.Log("deleting remote cluster")
				if err := c.DeleteRemote(ccr.RemoteName, true); err != nil {
					t.Logf("failed to delete remote cluster:%v", err)
				}
				t.Log("cleaning completed")
			},
			WantErr: true,
		},
		{
			Name:    "delete autofollow rule",
			Wrapper: testWrapper(),
			CaseInput: replication.DeleteAutofollowReq{
				Header: nil,
				Body: replication.DeleteAutofollowBody{
					Name:        autoFollowRequest.Body.Name,
					LeaderAlias: getCCR().RemoteName,
				},
			},
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(getCCR(), true), "expected to configure remote cluster")
				t.Log("creating autofollow rule")
				assert.NoError(t, c.CreateAutofollowRule(autoFollowRequest, true), "expected to create autofollow rule")
				t.Log("pre-configuration completed")
			},
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				if err := c.DeleteRemote(getCCR().RemoteName, true); err != nil {
					t.Logf("failed to delete remote cluster:%v", err)
				}
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Cleanup(func() {
				if tt.PostFunc != nil {
					tt.PostFunc(t, tt.Wrapper)
				}
			})

			if tt.ConfigureFunc != nil {
				log.Info().Msg("executing configure func")
				tt.ConfigureFunc(t, tt.Wrapper)
			}

			executionErr := tt.Wrapper.DeleteAutofollow(tt.CaseInput.(replication.DeleteAutofollowReq), true)
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}
