package tagger

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable tagged diff report to w.
func WriteText(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}
	for _, r := range results {
		tags := joinTags(r.Tags)
		fmt.Fprintf(w, "[%s] %s  tags:%s\n", r.Diff.Kind, r.Diff.Key, tags)
		if r.Diff.ValueA != "" {
			fmt.Fprintf(w, "  A: %s\n", r.Diff.ValueA)
		}
		if r.Diff.ValueB != "" {
			fmt.Fprintf(w, "  B: %s\n", r.Diff.ValueB)
		}
	}
}

// WriteJSON writes a JSON-encoded tagged diff report to w.
func WriteJSON(w io.Writer, results []Result) error {
	type jsonEntry struct {
		Key    string   `json:"key"`
		Kind   string   `json:"kind"`
		ValueA string   `json:"value_a,omitempty"`
		ValueB string   `json:"value_b,omitempty"`
		Tags   []string `json:"tags"`
	}

	entries := make([]jsonEntry, 0, len(results))
	for _, r := range results {
		tags := make([]string, len(r.Tags))
		for i, t := range r.Tags {
			tags[i] = string(t)
		}
		entries = append(entries, jsonEntry{
			Key:    r.Diff.Key,
			Kind:   string(r.Diff.Kind),
			ValueA: r.Diff.ValueA,
			ValueB: r.Diff.ValueB,
			Tags:   tags,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func joinTags(tags []Tag) string {
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = string(t)
	}
	return strings.Join(parts, ",")
}
