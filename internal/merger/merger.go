// Package merger provides functionality to merge multiple .env files into
// a single unified map, with configurable conflict resolution strategies.
package merger

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how key conflicts are resolved when merging.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if any key is defined in more than one file.
	StrategyError
)

// ParseStrategy converts a string label to a Strategy constant.
func ParseStrategy(s string) (Strategy, error) {
	switch s {
	case "first":
		return StrategyFirst, nil
	case "last":
		return StrategyLast, nil
	case "error":
		return StrategyError, nil
	default:
		return StrategyFirst, fmt.Errorf("unknown merge strategy %q: must be first, last, or error", s)
	}
}

// Result holds the merged environment and metadata about the merge.
type Result struct {
	// Env is the merged key-value map.
	Env parser.EnvFile
	// Sources maps each key to the file path it was taken from.
	Sources map[string]string
}

// Keys returns a sorted-order-independent slice of all keys in the merged result.
// It is a convenience method for iterating over Env without importing maps or sort.
func (r *Result) Keys() []string {
	keys := make([]string, 0, len(r.Env))
	for k := range r.Env {
		keys = append(keys, k)
	}
	return keys
}

// Merge combines multiple named env files into one according to the given strategy.
// files is a slice of (label, EnvFile) pairs; label is used in error messages and
// source tracking.
func Merge(files []NamedEnv, strategy Strategy) (*Result, error) {
	result := &Result{
		Env:     make(parser.EnvFile),
		Sources: make(map[string]string),
	}

	for _, nf := range files {
		for key, value := range nf.Env {
			existing, exists := result.Env[key]
			_ = existing
			if !exists {
				result.Env[key] = value
				result.Sources[key] = nf.Label
				continue
			}
			switch strategy {
			case StrategyFirst:
				// keep existing, do nothing
			case StrategyLast:
				result.Env[key] = value
				result.Sources[key] = nf.Label
			case StrategyError:
				return nil, fmt.Errorf("conflict: key %q defined in both %q and %q",
					key, result.Sources[key], nf.Label)
			}
		}
	}
	return result, nil
}

// NamedEnv pairs an EnvFile with a human-readable label (e.g. a file path).
type NamedEnv struct {
	Label string
	Env   parser.EnvFile
}
