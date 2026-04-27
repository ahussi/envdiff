package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options configures report output.
type Options struct {
	Format  Format
	Colorize bool
}

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
)

// WriteText writes a human-readable diff report to w.
func WriteText(w io.Writer, results []diff.Result, fileA, fileB string, colorize bool) {
	if len(results) == 0 {
		fmt.Fprintf(w, "No differences found between %s and %s\n", fileA, fileB)
		return
	}

	fmt.Fprintf(w, "Comparing %s vs %s\n", fileA, fileB)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, r := range results {
		switch r.Kind {
		case diff.MissingInB:
			prefix := ""
			if colorize {
				prefix = colorCyan
			}
			fmt.Fprintf(w, "%s[MISSING in %s]%s KEY=%q VALUE=%q\n", prefix, fileB, colorReset, r.Key, r.ValueA)
		case diff.MissingInA:
			prefix := ""
			if colorize {
				prefix = colorYellow
			}
			fmt.Fprintf(w, "%s[MISSING in %s]%s KEY=%q VALUE=%q\n", prefix, fileA, colorReset, r.Key, r.ValueB)
		case diff.Mismatch:
			prefix := ""
			if colorize {
				prefix = colorRed
			}
			fmt.Fprintf(w, "%s[MISMATCH]%s KEY=%q %s=%q %s=%q\n", prefix, colorReset, r.Key, fileA, r.ValueA, fileB, r.ValueB)
		}
	}
}
