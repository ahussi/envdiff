// Package grouper groups diff results by a common key prefix (e.g. "DB_", "AWS_").
// Keys without an underscore are placed in a special "_ungrouped" bucket.
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Group represents a named collection of diff results that share a prefix.
type Group struct {
	Prefix  string
	Results []diff.Result
}

// Analyse partitions results by their first underscore-delimited prefix.
// Results whose keys contain no underscore are collected under "_ungrouped".
func Analyse(results []diff.Result) []Group {
	buckets := make(map[string][]diff.Result)

	for _, r := range results {
		prefix := extractPrefix(r.Key)
		buckets[prefix] = append(buckets[prefix], r)
	}

	groups := make([]Group, 0, len(buckets))
	for prefix, items := range buckets {
		groups = append(groups, Group{
			Prefix:  prefix,
			Results: items,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		// _ungrouped always last
		if groups[i].Prefix == "_ungrouped" {
			return false
		}
		if groups[j].Prefix == "_ungrouped" {
			return true
		}
		return groups[i].Prefix < groups[j].Prefix
	})

	return groups
}

// extractPrefix returns the portion of key before the first "_".
// If no underscore exists the constant "_ungrouped" is returned.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return "_ungrouped"
}
