package pinner

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable pin report to w.
func WriteText(w io.Writer, r Result) {
	if len(r.Pinned) == 0 && len(r.Floating) == 0 {
		fmt.Fprintln(w, "No keys found.")
		return
	}

	if len(r.Pinned) > 0 {
		fmt.Fprintf(w, "Pinned keys (%d):\n", len(r.Pinned))
		for _, pk := range r.Pinned {
			fmt.Fprintf(w, "  %-30s = %s  [%s]\n", pk.Key, pk.Value, strings.Join(pk.Envs, ", "))
		}
	}

	if len(r.Floating) > 0 {
		fmt.Fprintf(w, "Floating keys (%d):\n", len(r.Floating))
		for _, k := range r.Floating {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
}

// jsonResult is the JSON-serialisable form of Result.
type jsonResult struct {
	Pinned   []jsonPinnedKey `json:"pinned"`
	Floating []string        `json:"floating"`
}

type jsonPinnedKey struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Envs  []string `json:"envs"`
}

// WriteJSON writes the pin report as a JSON object to w.
func WriteJSON(w io.Writer, r Result) error {
	out := jsonResult{
		Floating: r.Floating,
	}
	if out.Floating == nil {
		out.Floating = []string{}
	}
	for _, pk := range r.Pinned {
		out.Pinned = append(out.Pinned, jsonPinnedKey{
			Key:   pk.Key,
			Value: pk.Value,
			Envs:  pk.Envs,
		})
	}
	if out.Pinned == nil {
		out.Pinned = []jsonPinnedKey{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
