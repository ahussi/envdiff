package snapshotter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable drift report to w.
func WriteText(w io.Writer, snap *Snapshot, drifts []DriftEntry) {
	fmt.Fprintf(w, "Snapshot: %s (captured %s)\n", snap.Label, snap.CapturedAt.Format("2006-01-02 15:04:05 UTC"))

	if len(drifts) == 0 {
		fmt.Fprintln(w, "No drift detected.")
		return
	}

	sorted := make([]DriftEntry, len(drifts))
	copy(sorted, drifts)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "%d drift(s) detected:\n\n", len(sorted))
	for _, d := range sorted {
		switch d.Kind {
		case DriftAdded:
			fmt.Fprintf(w, "  [added]   %s = %q\n", d.Key, d.NewValue)
		case DriftRemoved:
			fmt.Fprintf(w, "  [removed] %s (was %q)\n", d.Key, d.OldValue)
		case DriftChanged:
			fmt.Fprintf(w, "  [changed] %s: %q -> %q\n", d.Key, d.OldValue, d.NewValue)
		}
	}
}

// WriteJSON writes the drift report as a JSON object to w.
func WriteJSON(w io.Writer, snap *Snapshot, drifts []DriftEntry) error {
	payload := struct {
		Label      string      `json:"label"`
		CapturedAt string      `json:"captured_at"`
		DriftCount int         `json:"drift_count"`
		Drifts     []DriftEntry `json:"drifts"`
	}{
		Label:      snap.Label,
		CapturedAt: snap.CapturedAt.Format("2006-01-02T15:04:05Z"),
		DriftCount: len(drifts),
		Drifts:     drifts,
	}
	if payload.Drifts == nil {
		payload.Drifts = []DriftEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
