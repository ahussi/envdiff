// Package ignorer provides functionality to skip keys matching
// user-supplied glob patterns when comparing .env files.
package ignorer

import (
	"path"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options holds the configuration for the ignorer.
type Options struct {
	// Patterns is a list of glob patterns (e.g. "*_SECRET", "AWS_*").
	Patterns []string
}

// Apply filters out any diff.Result whose key matches at least one of the
// provided glob patterns. Results that do not match any pattern are returned
// unchanged. An empty Patterns slice is a no-op.
func Apply(results []diff.Result, opts Options) ([]diff.Result, error) {
	if len(opts.Patterns) == 0 {
		return results, nil
	}

	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		matched, err := matchesAny(r.Key, opts.Patterns)
		if err != nil {
			return nil, err
		}
		if !matched {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

// matchesAny reports whether key matches any of the given glob patterns.
// Pattern matching is case-insensitive and uses path.Match semantics.
func matchesAny(key string, patterns []string) (bool, error) {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		matched, err := path.Match(strings.ToUpper(p), upper)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}
