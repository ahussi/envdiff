// Package snapshotter captures a point-in-time snapshot of a parsed .env file
// and compares it against a current state to detect drift over time.
package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file.
type Snapshot struct {
	Label     string            `json:"label"`
	CapturedAt time.Time        `json:"captured_at"`
	Entries   map[string]string `json:"entries"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, label string, entries map[string]string) error {
	snap := Snapshot{
		Label:      label,
		CapturedAt: time.Now().UTC(),
		Entries:    entries,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshotter: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshotter: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: read %s: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshotter: unmarshal: %w", err)
	}
	return &snap, nil
}

// DriftKind describes the type of change detected between snapshot and current.
type DriftKind string

const (
	DriftAdded   DriftKind = "added"   // key present in current but not snapshot
	DriftRemoved DriftKind = "removed" // key present in snapshot but not current
	DriftChanged DriftKind = "changed" // key present in both but value differs
)

// DriftEntry represents a single detected change.
type DriftEntry struct {
	Key      string    `json:"key"`
	Kind     DriftKind `json:"kind"`
	OldValue string    `json:"old_value,omitempty"`
	NewValue string    `json:"new_value,omitempty"`
}

// Diff compares a snapshot against current entries and returns detected drift.
func Diff(snap *Snapshot, current map[string]string) []DriftEntry {
	var results []DriftEntry

	for k, oldVal := range snap.Entries {
		newVal, ok := current[k]
		if !ok {
			results = append(results, DriftEntry{Key: k, Kind: DriftRemoved, OldValue: oldVal})
		} else if newVal != oldVal {
			results = append(results, DriftEntry{Key: k, Kind: DriftChanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, newVal := range current {
		if _, ok := snap.Entries[k]; !ok {
			results = append(results, DriftEntry{Key: k, Kind: DriftAdded, NewValue: newVal})
		}
	}

	return results
}
