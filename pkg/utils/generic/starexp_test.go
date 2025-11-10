package generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestWildcard validates matching functionality for strings against different patterns containing wildcards.
func TestWildcard(t *testing.T) {
	tests := []struct {
		name          string
		stringPattern string
		valueMap      map[string]bool // key string itself, value -> expected result of matching to the pattern
	}{
		{
			name:          "empty string pattern",
			stringPattern: "",
			valueMap: map[string]bool{
				"":          true,
				"something": false,
			},
		},
		{
			name:          "no wildcard",
			stringPattern: "aaa",
			valueMap: map[string]bool{
				"aaa":                  true,
				"some negative string": false,
			},
		},
		{
			name:          "wildcard only",
			stringPattern: "*",
			valueMap: map[string]bool{
				"aaa":                             true,
				"aaaa":                            true,
				"aaa this is the match bbb":       true,
				"aaabbb":                          true,
				"aaa1bbb":                         true,
				"aaa, some positive long string?": true,
				"is this tricky string? aaa":      true,
				"is this tricky string? aaabbb":   true,
			},
		},
		{
			name:          "wildcard in the beginning",
			stringPattern: "*aaa",
			valueMap: map[string]bool{
				"aaa":                        true,
				"aaaa":                       true,
				"is this tricky string? aaa": true,
				"some negative string":       false,
			},
		},
		{
			name:          "wildcard in the end",
			stringPattern: "aaa*",
			valueMap: map[string]bool{
				"aaa":                            true,
				"aaaa":                           true,
				"aaa, some positive long string": true,
				"is this tricky string? aaa":     false,
			},
		},
		{
			name:          "wildcard in the middle",
			stringPattern: "aaa*bbb",
			valueMap: map[string]bool{
				"aaa":                             false,
				"aaaa":                            false,
				"aaa this is the match bbb":       true,
				"aaabbb":                          true,
				"aaa1bbb":                         true,
				"aaa, some positive long string?": false,
				"is this tricky string? aaa":      false,
				"is this tricky string? aaabbb":   false,
			},
		},
		{
			name:          "wildcard from both sides",
			stringPattern: "*a*",
			valueMap: map[string]bool{
				"aaa":                             true,
				"aaaa":                            true,
				"aaa this is the match bbb":       true,
				"aaabbb":                          true,
				"aaa1bbb":                         true,
				"aaa, some positive long string?": true,
				"is this tricky string? aaa":      true,
				"is this tricky string? aaabbb":   true,
			},
		},
		{
			name:          "double wildcards",
			stringPattern: "**",
			valueMap: map[string]bool{
				"aaa":                             true,
				"aaaa":                            true,
				"aaa this is the match bbb":       true,
				"aaabbb":                          true,
				"aaa1bbb":                         true,
				"aaa, some positive long string?": true,
				"is this tricky string? aaa":      true,
				"is this tricky string? aaabbb":   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fx := GetMatchFunc(tt.stringPattern)
			for testString, expectedMatchResult := range tt.valueMap {
				assert.Equal(t, expectedMatchResult, fx(testString), "function output doesn't match with expected result")
			}
		})
	}
}
