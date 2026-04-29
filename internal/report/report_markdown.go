package report

import (
	"fmt"
	"io"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/formatter"
)

// WriteMarkdown writes a Markdown-formatted diff report to w.
func WriteMarkdown(w io.Writer, results []diff.Result, fileA, fileB string) error {
	fmt.Fprintf(w, "# envdiff report\n\n")
	fmt.Fprintf(w, "Comparing `%s` → `%s`\n\n", fileA, fileB)

	if len(results) == 0 {
		fmt.Fprintln(w, "✅ No differences found.")
		return nil
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

	fmt.Fprintf(w, "\n_%d issue(s) found._\n", len(results))
	return nil
}

func escapeMarkdown(s string) string {
	if s == "" {
		return "_(empty)_"
	}
	return "`" + s + "`"
}
