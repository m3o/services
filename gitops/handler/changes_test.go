package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineChanges(t *testing.T) {
	tt := []struct {
		Name    string
		Commits []commit
		Result  map[string]changeType
	}{
		{
			Name: "ServiceCreated",
			Commits: []commit{
				{
					Added: []string{
						"foo/main.go",
						"foo/handler/handler.go",
						"foo/api/main.go",
						"foo/api/handler/handler.go",
					},
				},
			},
			Result: map[string]changeType{
				"foo":     created,
				"foo/api": created,
			},
		},
		{
			Name: "ServiceDeleted",
			Commits: []commit{
				{
					Removed: []string{
						"foo/main.go",
						"foo/handler/handler.go",
						"foo/api/main.go",
						"foo/api/handler/handler.go",
					},
				},
			},
			Result: map[string]changeType{
				"foo":     deleted,
				"foo/api": deleted,
			},
		},
		{
			Name: "ServiceModified",
			Commits: []commit{
				{
					Modified: []string{
						"foo/main.go",
						"bar/handler/handler.go",
					},
				},
			},
			Result: map[string]changeType{},
		},
		{
			Name: "MultipleServices",
			Commits: []commit{
				{
					Added: []string{
						"foo/main.go",
						"foo/api/main.go",
					},
					Modified: []string{
						"bar/main.go",
						"bar/handler/handler.go",
					},
					Removed: []string{"bar/api/main.go"},
				},
			},
			Result: map[string]changeType{
				"foo":     created,
				"foo/api": created,
				"bar/api": deleted,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			res := determineChanges(tc.Commits)
			assert.Equal(t, tc.Result, res, "Expected the results to match")
		})
	}
}
