package anchorer

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable anchor report to w.
func WriteText(w io.Writer, r Result) {
	if len(r.Anchors) == 0 {
		fmt.Fprintln(w, "Anchors: none")
	} else {
		fmt.Fprintf(w, "Anchors (%d):\n", len(r.Anchors))
		for _, k := range r.Anchors {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}

	if len(r.Unique) == 0 {
		fmt.Fprintln(w, "Unique keys: none")
		return
	}

	labels := make([]string, 0, len(r.Unique))
	for l := range r.Unique {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	for _, label := range labels {
		keys := r.Unique[label]
		if len(keys) == 0 {
			continue
		}
		fmt.Fprintf(w, "Unique to %s (%d):\n", label, len(keys))
		for _, k := range keys {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
}

// jsonResult is the JSON-serialisable form of Result.
type jsonResult struct {
	Anchors []string            `json:"anchors"`
	Unique  map[string][]string `json:"unique"`
}

// WriteJSON writes a JSON-encoded anchor report to w.
func WriteJSON(w io.Writer, r Result) error {
	out := jsonResult{
		Anchors: r.Anchors,
		Unique:  r.Unique,
	}
	if out.Anchors == nil {
		out.Anchors = []string{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
