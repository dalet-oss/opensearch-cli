package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpensearchWrapper_DeleteIndex(t *testing.T) {
	indexToDelete := "tc-create-index-to-delete"
	tests := []TestCase{
		{
			Name:          "delete index that doesn't exist",
			Wrapper:       testWrapper(),
			CaseInput:     "not-existing-index",
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
		},
		{
			Name:      "delete index",
			Wrapper:   testWrapper(),
			CaseInput: indexToDelete,
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Logf("creating index %s", indexToDelete)
				assert.NoError(t, c.CreateIndex(indexToDelete), "expected to create index")
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
			executionErr := tt.Wrapper.DeleteIndex(tt.CaseInput.(string))
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}
func TestOpensearchWrapper_CreateIndex(t *testing.T) {
	tests := []TestCase{
		{
			Name:          "create index",
			Wrapper:       testWrapper(),
			CaseInput:     "tc-create-index",
			ConfigureFunc: nil,
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("Cleaning things")
				assert.NoError(t, c.DeleteIndex("tc-create-index"), "expected to delete index")
			},
			WantErr: false,
		},
		{
			Name:      "create index that already exists",
			Wrapper:   testWrapper(),
			CaseInput: "tc-create-index-fail-already-exists",
			ConfigureFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("creating index")
				assert.NoError(t, c.CreateIndex("tc-create-index-fail-already-exists"), "expected to create index")
			},
			PostFunc: func(t *testing.T, c *OpensearchWrapper) {
				t.Log("Cleaning things")
				assert.NoError(t, c.DeleteIndex("tc-create-index-fail-already-exists"), "expected to delete index")
			},
			WantErr: true,
		},
		{
			Name:          "create index with invalid name",
			Wrapper:       testWrapper(),
			CaseInput:     "",
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       true,
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
			executionErr := tt.Wrapper.CreateIndex(tt.CaseInput.(string))
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}
func TestOpensearchWrapper_GetIndexList(t *testing.T) {
	tests := []TestCase{
		{
			Name:          "get index list",
			Wrapper:       testWrapper(),
			CaseInput:     nil,
			ConfigureFunc: nil,
			PostFunc:      nil,
			WantErr:       false,
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
			_, executionErr := tt.Wrapper.GetIndexList()
			if tt.WantErr {
				assert.Error(t, executionErr, "expected to get error")
			} else {
				assert.NoError(t, executionErr, "expected to get no error")
			}
		})
	}
}
