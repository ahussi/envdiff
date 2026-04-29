package filter

import "github.com/yourusername/envdiff/internal/diff"

// Options holds filtering criteria for diff results.
type Options struct {
	// OnlyKinds restricts results to specific diff kinds.
	// Valid values: "missing_in_b", "missing_in_a", "mismatch"
	OnlyKinds []string

	// KeyPrefix filters results to keys starting with the given prefix.
	KeyPrefix string
}

// Apply returns a filtered subset of the given diff results based on Options.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if len(opts.OnlyKinds) == 0 && opts.KeyPrefix == "" {
		return results
	}

	kindSet := make(map[string]bool, len(opts.OnlyKinds))
	for _, k := range opts.OnlyKinds {
		kindSet[k] = true
	}

	var filtered []diff.Result
	for _, r := range results {
		if opts.KeyPrefix != "" && !hasPrefix(r.Key, opts.KeyPrefix) {
			continue
		}
		if len(kindSet) > 0 && !kindSet[kindLabel(r.Kind)] {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

func hasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

func kindLabel(k diff.Kind) string {
	switch k {
	case diff.MissingInB:
		return "missing_in_b"
	case diff.MissingInA:
		return "missing_in_a"
	case diff.Mismatch:
		return "mismatch"
	default:
		return ""
	}
}
