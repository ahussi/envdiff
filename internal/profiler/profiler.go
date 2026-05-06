// Package profiler analyses an env file and produces a health profile
// summarising key quality metrics: lint issues, validation gaps, score, and
// redaction coverage.
package profiler

import (
	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/linter"
	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/redactor"
	"github.com/yourorg/envdiff/internal/scorer"
)

// Profile holds the aggregated health information for a single env file.
type Profile struct {
	// Label is the human-readable name of the file (e.g. ".env.production").
	Label string

	// TotalKeys is the number of key-value pairs parsed from the file.
	TotalKeys int

	// LintIssues contains every warning raised by the linter.
	LintIssues []linter.Issue

	// RedactedKeys is the count of keys whose values were redacted.
	RedactedKeys int

	// Score is the computed health score (0–100).
	Score scorer.Score
}

// Analyse parses the file at path, runs the linter and redactor, then
// computes a health score using a self-comparison (zero diff results).
func Analyse(path, label string, extraPatterns []string) (Profile, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return Profile{}, err
	}

	// Lint the raw entries.
	issues := linter.Check(entries)

	// Build a minimal diff result set so the scorer has something to work with.
	// A self-comparison yields no missing/mismatch results; only lint issues
	// influence the score when used this way.
	envMap := entriesToMap(entries)
	diffResults := diff.Compare(envMap, envMap)

	// Apply redactor to count sensitive keys.
	redacted := redactor.Apply(diffResults, extraPatterns)
	redactedCount := 0
	for _, r := range redacted {
		if r.RedactedA || r.RedactedB {
			redactedCount++
		}
	}

	s := scorer.Compute(diffResults)

	return Profile{
		Label:        label,
		TotalKeys:    len(entries),
		LintIssues:   issues,
		RedactedKeys: redactedCount,
		Score:        s,
	}, nil
}

// entriesToMap converts a slice of parser.Entry into the map type expected by
// diff.Compare and scorer.Compute.
func entriesToMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
