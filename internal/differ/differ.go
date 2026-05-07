// Package differ provides utilities for computing a similarity score
// between two env maps, expressed as a percentage of matching keys.
package differ

import (
	"fmt"
	"strings"
)

// Result holds the overlap analysis between two env maps.
type Result struct {
	// TotalUnique is the number of distinct keys across both maps.
	TotalUnique int
	// Shared is the number of keys present in both maps.
	Shared int
	// OnlyInA contains keys present only in the first map.
	OnlyInA []string
	// OnlyInB contains keys present only in the second map.
	OnlyInB []string
	// SimilarityPct is Shared/TotalUnique expressed as a percentage (0–100).
	SimilarityPct float64
}

// Analyse computes the overlap between two env maps keyed by variable name.
func Analyse(a, b map[string]string) Result {
	union := make(map[string]struct{})
	for k := range a {
		union[k] = struct{}{}
	}
	for k := range b {
		union[k] = struct{}{}
	}

	var shared int
	var onlyA, onlyB []string

	for k := range union {
		_, inA := a[k]
		_, inB := b[k]
		switch {
		case inA && inB:
			shared++
		case inA:
			onlyA = append(onlyA, k)
		default:
			onlyB = append(onlyB, k)
		}
	}

	total := len(union)
	var pct float64
	if total > 0 {
		pct = float64(shared) / float64(total) * 100
	}

	return Result{
		TotalUnique:   total,
		Shared:        shared,
		OnlyInA:       onlyA,
		OnlyInB:       onlyB,
		SimilarityPct: pct,
	}
}

// Summary returns a human-readable one-line description of the result.
func Summary(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "similarity: %.1f%% (%d/%d shared)",
		r.SimilarityPct, r.Shared, r.TotalUnique)
	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(&sb, ", only-in-A: %d", len(r.OnlyInA))
	}
	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(&sb, ", only-in-B: %d", len(r.OnlyInB))
	}
	return sb.String()
}
