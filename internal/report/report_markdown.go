package report

import (
	"fmt"
	"io"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/formatter"
)

// WriteMarkdown writes a Markdown-formatted diff report to w.
// It renders a summary header, a comparison table of all differing keys,
// and a footer with the total issue count.
func WriteMarkdown(w io.Writer, results []diff.Result, fileA, fileB string) error {
	if _, err := fmt.Fprintf(w, "# envdiff report\n\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Comparing `%s` → `%s`\n\n", fileA, fileB); err != nil {
		return err
	}

	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "✅ No differences found.")
		return err
	}

	fmt.Fprintf(w, "| Key | Kind | Value A | Value B |\n")
	fmt.Fprintf(w, "|-----|------|---------|---------|\n")

	for _, r := range results {
		kind := formatter.KindLabel(kindString(r.Kind), formatter.StyleMarkdown)
		key := formatter.KeyLabel(r.Key, formatter.StyleMarkdown)
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n",
			key,
			kind,
			escapeMarkdown(r.ValueA),
			escapeMarkdown(r.ValueB),
		)
	}

	_, err := fmt.Fprintf(w, "\n_%d issue(s) found._\n", len(results))
	return err
}

func escapeMarkdown(s string) string {
	if s == "" {
		return "_(empty)_"
	}
	return "`" + s + "`"
}
