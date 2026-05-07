// Package renamer provides utilities for detecting and suggesting key renames
// between two .env files. A rename is inferred when a key disappears in one
// environment and a new key with a matching value appears in the other.
package renamer

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Suggestion represents a probable key rename between two environments.
type Suggestion struct {
	OldKey string
	NewKey string
	Value  string
}

// String returns a human-readable description of the rename suggestion.
func (s Suggestion) String() string {
	return fmt.Sprintf("%s -> %s (value: %q)", s.OldKey, s.NewKey, s.Value)
}

// Detect analyses a slice of diff results and returns rename suggestions.
// A rename is inferred when a key missing in B shares its value with a key
// missing in A — implying the key was renamed rather than removed/added.
func Detect(results []diff.Result) []Suggestion {
	// Index keys missing in B by their value (present only in A).
	missingInB := make(map[string]string) // value -> oldKey
	for _, r := range results {
		if r.Kind == diff.MissingInB {
			if r.ValueA != "" {
				missingInB[r.ValueA] = r.Key
			}
		}
	}

	var suggestions []Suggestion
	for _, r := range results {
		if r.Kind == diff.MissingInA {
			if oldKey, ok := missingInB[r.ValueB]; ok {
				suggestions = append(suggestions, Suggestion{
					OldKey: oldKey,
					NewKey: r.Key,
					Value:  r.ValueB,
				})
			}
		}
	}
	return suggestions
}
