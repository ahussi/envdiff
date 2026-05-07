package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/snapshotter"
)

func makeEntries() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entries := makeEntries()
	if err := snapshotter.Save(path, "prod", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := snapshotter.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if snap.Label != "prod" {
		t.Errorf("label: got %q, want %q", snap.Label, "prod")
	}
	if snap.Entries["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q", snap.Entries["DB_HOST"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshotter.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := snapshotter.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiff_NoChanges(t *testing.T) {
	snap := &snapshotter.Snapshot{Label: "test", CapturedAt: time.Now(), Entries: makeEntries()}
	drifts := snapshotter.Diff(snap, makeEntries())
	if len(drifts) != 0 {
		t.Errorf("expected no drift, got %d", len(drifts))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	snap := &snapshotter.Snapshot{Label: "test", CapturedAt: time.Now(), Entries: makeEntries()}
	current := makeEntries()
	current["NEW_KEY"] = "value"

	drifts := snapshotter.Diff(snap, current)
	if len(drifts) != 1 || drifts[0].Kind != snapshotter.DriftAdded {
		t.Errorf("expected 1 added drift, got %+v", drifts)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	snap := &snapshotter.Snapshot{Label: "test", CapturedAt: time.Now(), Entries: makeEntries()}
	current := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	drifts := snapshotter.Diff(snap, current)
	if len(drifts) != 1 || drifts[0].Kind != snapshotter.DriftRemoved || drifts[0].Key != "APP_ENV" {
		t.Errorf("expected 1 removed drift, got %+v", drifts)
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	snap := &snapshotter.Snapshot{Label: "test", CapturedAt: time.Now(), Entries: makeEntries()}
	current := makeEntries()
	current["DB_HOST"] = "remotehost"

	drifts := snapshotter.Diff(snap, current)
	if len(drifts) != 1 || drifts[0].Kind != snapshotter.DriftChanged {
		t.Errorf("expected 1 changed drift, got %+v", drifts)
	}
	if drifts[0].OldValue != "localhost" || drifts[0].NewValue != "remotehost" {
		t.Errorf("unexpected values: %+v", drifts[0])
	}
}
