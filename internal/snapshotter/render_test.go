package snapshotter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/snapshotter"
)

func makeSnap() *snapshotter.Snapshot {
	return &snapshotter.Snapshot{
		Label:      "staging",
		CapturedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Entries:    makeEntries(),
	}
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	snapshotter.WriteText(&buf, makeSnap(), nil)

	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Error("expected label in output")
	}
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	drifts := []snapshotter.DriftEntry{
		{Key: "DB_HOST", Kind: snapshotter.DriftChanged, OldValue: "localhost", NewValue: "remotehost"},
		{Key: "NEW_KEY", Kind: snapshotter.DriftAdded, NewValue: "val"},
		{Key: "OLD_KEY", Kind: snapshotter.DriftRemoved, OldValue: "gone"},
	}

	var buf bytes.Buffer
	snapshotter.WriteText(&buf, makeSnap(), drifts)
	out := buf.String()

	for _, want := range []string{"[added]", "[removed]", "[changed]", "DB_HOST", "NEW_KEY", "OLD_KEY"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	drifts := []snapshotter.DriftEntry{
		{Key: "X", Kind: snapshotter.DriftAdded, NewValue: "1"},
	}

	var buf bytes.Buffer
	if err := snapshotter.WriteJSON(&buf, makeSnap(), drifts); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	out := buf.String()

	for _, want := range []string{`"label"`, `"staging"`, `"drift_count"`, `"drifts"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON:\n%s", want, out)
		}
	}
}

func TestWriteJSON_EmptyDrifts(t *testing.T) {
	var buf bytes.Buffer
	if err := snapshotter.WriteJSON(&buf, makeSnap(), nil); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drift_count": 0`) {
		t.Errorf("expected drift_count 0, got:\n%s", out)
	}
	if !strings.Contains(out, `"drifts": []`) {
		t.Errorf("expected empty drifts array, got:\n%s", out)
	}
}
