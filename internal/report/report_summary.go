package report

import (
	"io"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

// WriteSummaryText appends a plain-text summary line after the diff report.
func WriteSummaryText(w io.Writer, results []diff.Result) error {
	if err := WriteText(w, results); err != nil {
		return err
	}
	s := summary.Compute(results)
	summary.WriteText(w, s)
	return nil
}

// WriteSummaryJSON writes a JSON object containing both the diff results
// and a top-level summary field.
func WriteSummaryJSON(w io.Writer, results []diff.Result) error {
	// Delegate full diff JSON then append summary on its own line for
	// machine-readable consumers that process line-delimited JSON.
	if err := WriteJSON(w, results); err != nil {
		return err
	}
	s := summary.Compute(results)
	summary.WriteJSON(w, s)
	return nil
}
