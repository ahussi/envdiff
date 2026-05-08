// Package tagger assigns semantic tags to diff results based on key naming
// conventions and value patterns, enabling richer categorisation of entries.
package tagger

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Tag represents a semantic label attached to a diff result.
type Tag string

const (
	TagSecret  Tag = "secret"
	TagURL     Tag = "url"
	TagFeature Tag = "feature_flag"
	TagPort    Tag = "port"
	TagPath    Tag = "path"
	TagUnknown Tag = "unknown"
)

// Result wraps a diff.Result with one or more tags.
type Result struct {
	Diff diff.Result
	Tags []Tag
}

// Annotate applies tags to each diff result and returns the annotated slice.
func Annotate(results []diff.Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Diff: r,
			Tags: tagsForKey(r.Key),
		})
	}
	return out
}

func tagsForKey(key string) []Tag {
	upper := strings.ToUpper(key)
	var tags []Tag

	if containsAny(upper, []string{"SECRET", "PASSWORD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL"}) {
		tags = append(tags, TagSecret)
	}
	if containsAny(upper, []string{"URL", "HOST", "ENDPOINT", "ADDR", "DSN"}) {
		tags = append(tags, TagURL)
	}
	if containsAny(upper, []string{"FEATURE", "FLAG", "ENABLE", "DISABLE", "TOGGLE"}) {
		tags = append(tags, TagFeature)
	}
	if containsAny(upper, []string{"PORT"}) {
		tags = append(tags, TagPort)
	}
	if containsAny(upper, []string{"PATH", "DIR", "DIRECTORY", "FILE"}) {
		tags = append(tags, TagPath)
	}
	if len(tags) == 0 {
		tags = append(tags, TagUnknown)
	}
	return tags
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
