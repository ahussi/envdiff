package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortBy defines the field to sort results by.
type SortBy string

const (
	SortByKey  SortBy = "key"
	SortByKind SortBy = "kind"
)

// Order defines the sort direction.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Options configures the sorting behaviour.
type Options struct {
	By    SortBy
	Order Order
}

// Apply sorts a slice of diff.Result according to the given options.
// If opts.By is empty, the original order is preserved.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if opts.By == "" || len(results) == 0 {
		return results
	}

	out := make([]diff.Result, len(results))
	copy(out, results)

	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch opts.By {
		case SortByKind:
			less = kindRank(out[i].Kind) < kindRank(out[j].Kind)
			if kindRank(out[i].Kind) == kindRank(out[j].Kind) {
				less = out[i].Key < out[j].Key
			}
		default: // SortByKey
			less = out[i].Key < out[j].Key
		}
		if opts.Order == OrderDesc {
			return !less
		}
		return less
	})

	return out
}

// kindRank assigns a numeric rank to each diff kind for stable ordering.
func kindRank(k diff.Kind) int {
	switch k {
	case diff.MissingInB:
		return 0
	case diff.MissingInA:
		return 1
	case diff.Mismatch:
		return 2
	default:
		return 3
	}
}
