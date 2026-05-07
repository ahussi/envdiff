package cascader

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes the cascade result as a human-readable table.
func WriteText(w io.Writer, r *Result) {
	if len(r.Resolved) == 0 {
		fmt.Fprintln(w, "No keys resolved.")
		return
	}

	keys := sortedKeys(r.Resolved)
	fmt.Fprintf(w, "%-40s %-20s %s\n", "KEY", "SOURCE", "VALUE")
	fmt.Fprintf(w, "%-40s %-20s %s\n", "---", "------", "-----")
	for _, k := range keys {
		source := r.Provenance[k]
		val := r.Resolved[k]
		if len(val) > 40 {
			val = val[:37] + "..."
		}
		fmt.Fprintf(w, "%-40s %-20s %s\n", k, source, val)
	}
}

type jsonResult struct {
	Layers     []string          `json:"layers"`
	Resolved   map[string]string `json:"resolved"`
	Provenance map[string]string `json:"provenance"`
}

// WriteJSON writes the cascade result as JSON.
func WriteJSON(w io.Writer, r *Result) error {
	labels := make([]string, len(r.Layers))
	for i, l := range r.Layers {
		labels[i] = l.Label
	}
	out := jsonResult{
		Layers:     labels,
		Resolved:   r.Resolved,
		Provenance: r.Provenance,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
