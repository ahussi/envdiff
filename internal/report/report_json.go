package report

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// JSONEntry is the serialisable form of a single diff result.
type JSONEntry struct {
	Kind   string `json:"kind"`
	Key    string `json:"key"`
	ValueA string `json:"value_a,omitempty"`
	ValueB string `json:"value_b,omitempty"`
}

// JSONReport is the top-level JSON output structure.
type JSONReport struct {
	FileA   string      `json:"file_a"`
	FileB   string      `json:"file_b"`
	Count   int         `json:"diff_count"`
	Entries []JSONEntry `json:"diffs"`
}

// WriteJSON writes a JSON-encoded diff report to w.
func WriteJSON(w io.Writer, results []diff.Result, fileA, fileB string) error {
	entries := make([]JSONEntry, 0, len(results))
	for _, r := range results {
		entry := JSONEntry{
			Kind:   kindString(r.Kind),
			Key:    r.Key,
			ValueA: r.ValueA,
			ValueB: r.ValueB,
		}
		entries = append(entries, entry)
	}

	report := JSONReport{
		FileA:   fileA,
		FileB:   fileB,
		Count:   len(entries),
		Entries: entries,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func kindString(k diff.DiffKind) string {
	switch k {
	case diff.MissingInB:
		return "missing_in_b"
	case diff.MissingInA:
		return "missing_in_a"
	case diff.Mismatch:
		return "mismatch"
	default:
		return "unknown"
	}
}
