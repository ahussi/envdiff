// Package scorer computes a numeric health score for a set of diff results.
// The score ranges from 0 (worst) to 100 (perfect match).
package scorer

import (
	"fmt"
	"math"

	"github.com/user/envdiff/internal/diff"
)

// Weights assigned to each kind of difference.
const (
	weightMissing  = 3.0
	weightMismatch = 1.5
	weightExtra    = 1.0
)

// Result holds the computed score and a human-readable grade.
type Result struct {
	Score    float64 // 0–100
	Grade    string  // A, B, C, D, F
	Penalty  float64 // total weighted penalty
	Total    int     // total keys considered
}

// Compute derives a health score from the provided diff results.
// total is the number of unique keys across both files.
func Compute(results []diff.Result, total int) Result {
	if total <= 0 {
		return Result{Score: 100, Grade: "A", Penalty: 0, Total: 0}
	}

	var penalty float64
	for _, r := range results {
		switch r.Kind {
		case diff.MissingInB, diff.MissingInA:
			penalty += weightMissing
		case diff.Mismatch:
			penalty += weightMismatch
		default:
			penalty += weightExtra
		}
	}

	// Normalise: max possible penalty if every key were missing
	maxPenalty := float64(total) * weightMissing
	raw := 100.0 * (1.0 - penalty/maxPenalty)
	score := math.Max(0, math.Min(100, raw))

	return Result{
		Score:   math.Round(score*10) / 10,
		Grade:   grade(score),
		Penalty: penalty,
		Total:   total,
	}
}

// String returns a short summary of the result.
func (r Result) String() string {
	return fmt.Sprintf("score=%.1f grade=%s penalty=%.1f total=%d",
		r.Score, r.Grade, r.Penalty, r.Total)
}

func grade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
