package api

import (
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCCRCreateOpts_BuildCCRParams tests the BuildCCRParams method of the CCRCreateOpts struct
func TestCCRCreateOpts_BuildCCRParams(t *testing.T) {
	tests := []struct {
		Name string
		Opts CCRCreateOpts
		Want map[string]interface{}
	}{
		{
			Name: "test defaults",
			Opts: CCRCreateOpts{
				Type:       "",
				Mode:       "",
				RemoteName: "remote",
				RemoteAddr: "remote.fake:9300",
			},
			Want: map[string]interface{}{
				"persistent": map[string]interface{}{
					"cluster": map[string]interface{}{
						"remote": map[string]interface{}{
							"remote": map[string]interface{}{
								"mode":          "proxy",
								"proxy_address": "remote.fake:9300",
							},
						},
					},
				},
			},
		},
		{
			Name: "test diff type",
			Opts: CCRCreateOpts{
				Type:       "fake-type",
				Mode:       "",
				RemoteName: "remote",
				RemoteAddr: "remote.fake:9300",
			},
			Want: map[string]interface{}{
				"fake-type": map[string]interface{}{
					"cluster": map[string]interface{}{
						"remote": map[string]interface{}{
							"remote": map[string]interface{}{
								"mode":          "proxy",
								"proxy_address": "remote.fake:9300",
							},
						},
					},
				},
			},
		},
		{
			Name: "test diff mode",
			Opts: CCRCreateOpts{
				Type:       "fake-type",
				Mode:       "fake-mode",
				RemoteName: "remote",
				RemoteAddr: "remote.fake:9300",
			},
			Want: map[string]interface{}{
				"fake-type": map[string]interface{}{
					"cluster": map[string]interface{}{
						"remote": map[string]interface{}{
							"remote": map[string]interface{}{
								"mode":          "fake-mode",
								"proxy_address": "remote.fake:9300",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, string(printutils.MarshalJSONOrDie(tt.Want)), string(tt.Opts.BuildCCRParams()), "expected to get the same result")
		})
	}
}

// TestOpensearchWrapper_CreateRemoteCluster tests the CreateRemoteCluster method of the OpensearchWrapper struct
func TestOpensearchWrapper_ConfigureRemoteCluster(t *testing.T) {
	tests := []OSSingleContainerTest{
		{
			Name:    "Trying to configure a remote cluster with a non-existent type",
			Wrapper: testWrapper(),
			CaseInput: CCRCreateOpts{
				Type:       "fake-type",
				Mode:       "",
				RemoteName: "remote-fail-type",
				RemoteAddr: "remote.fake:9300",
			},
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:    "Trying to configure a remote cluster with a non-existent mode",
			Wrapper: testWrapper(),
			CaseInput: CCRCreateOpts{
				Type:       "",
				Mode:       "fake-mode",
				RemoteName: "remote-fail-mode",
				RemoteAddr: "remote.fake:9300",
			},
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:    "Normal remote config except server address",
			Wrapper: testWrapper(),
			CaseInput: CCRCreateOpts{
				Type:       "",
				Mode:       "",
				RemoteName: "remote-fail-server",
				RemoteAddr: "remote.fake:9300",
			},
			ConfigureFunc: nil,
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("cleaning func for: Normal remote config except server address")
				if err := c.DeleteRemote("remote-fail-server", true); err != nil {
					t.Logf("failed to delete remote cluster: %s ", err)
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
			executionErr := tt.Wrapper.ConfigureRemoteCluster(tt.CaseInput.(CCRCreateOpts), true)
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}

func TestOpensearchWrapper_DeleteRemote(t *testing.T) {
	const testRemoteName = "test-delete-remote"
	tests := []OSSingleContainerTest{
		{
			Name:          "delete remote cluster that doesn't exist",
			Wrapper:       testWrapper(),
			CaseInput:     "remote-doesn't-exist",
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:      "delete remote that does exist",
			Wrapper:   testWrapper(),
			CaseInput: testRemoteName,
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(CCRCreateOpts{
					Type:       "",
					Mode:       "",
					RemoteName: testRemoteName,
					RemoteAddr: "fake.local:9300",
				}, true), "expected to configure remote cluster")
				t.Log("pre-configuration completed")
			},
			PostFunc: nil,
			WantErr:  false,
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
			executionErr := tt.Wrapper.DeleteRemote(tt.CaseInput.(string), true)
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}

func TestOpensearchWrapper_GetRemoteSettings(t *testing.T) {
	remoteSettings := CCRCreateOpts{
		Type:       "",
		Mode:       "",
		RemoteName: "normal-remote",
		RemoteAddr: "fake.local:9300",
	}
	tests := []OSSingleContainerTest{
		{
			Name:          "get remote settings that doesn't exist",
			Wrapper:       testWrapper(),
			CaseInput:     "remote-doesn't-exist",
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:      "get remote settings that does exist",
			Wrapper:   testWrapper(),
			CaseInput: remoteSettings,
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				assert.NotNil(t, c)
				t.Log("configuring remote cluster")
				assert.NoError(t, c.ConfigureRemoteCluster(remoteSettings, true), "expected to configure remote cluster")
				t.Log("pre-configuration completed")
			},
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("cleaning func for: get remote settings that does exist")
				if err := c.DeleteRemote(remoteSettings.RemoteName, true); err != nil {
					t.Logf("failed to delete remote cluster: %s ", err)
				}
			},
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
			// we have to disable raw output because it's printing response as is w/o throwing an error
			executionErr := tt.Wrapper.GetRemoteSettings(false)
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}
