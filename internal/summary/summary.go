package summary

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Stats holds aggregate counts for a diff result set.
type Stats struct {
	Total    int
	Missing  int
	Extra    int
	Mismatch int
}

// Compute calculates statistics from a slice of diff results.
func Compute(results []diff.Result) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Kind {
		case diff.MissingInB:
			s.Missing++
		case diff.MissingInA:
			s.Extra++
		case diff.ValueMismatch:
			s.Mismatch++
		}
	}
	return s
}

// WriteText writes a human-readable summary line to w.
func WriteText(w io.Writer, s Stats) {
	if s.Total == 0 {
		fmt.Fprintln(w, "Summary: no differences found.")
		return
	}
	fmt.Fprintf(w, "Summary: %d difference(s) — %d missing, %d extra, %d mismatched\n",
		s.Total, s.Missing, s.Extra, s.Mismatch)
}

// WriteJSON writes a JSON summary object to w.
func WriteJSON(w io.Writer, s Stats) {
	fmt.Fprintf(w, `{"total":%d,"missing":%d,"extra":%d,"mismatch":%d}`,
		s.Total, s.Missing, s.Extra, s.Mismatch)
	fmt.Fprintln(w)
}
