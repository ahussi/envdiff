package anchorer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/anchorer"
)

func makeResult(anchors []string, unique map[string][]string) anchorer.Result {
	return anchorer.Result{Anchors: anchors, Unique: unique}
}

func TestWriteText_WithAnchors(t *testing.T) {
	r := makeResult([]string{"HOST", "PORT"}, nil)
	var buf bytes.Buffer
	anchorer.WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "HOST") || !strings.Contains(out, "PORT") {
		t.Errorf("expected anchor keys in output, got: %s", out)
	}
}

func TestWriteText_NoAnchors(t *testing.T) {
	r := makeResult(nil, nil)
	var buf bytes.Buffer
	anchorer.WriteText(&buf, r)
	if !strings.Contains(buf.String(), "none") {
		t.Errorf("expected 'none' in output for empty anchors")
	}
}

func TestWriteText_UniqueKeys(t *testing.T) {
	r := makeResult([]string{"SHARED"}, map[string][]string{
		"dev": {"DEV_ONLY"},
	})
	var buf bytes.Buffer
	anchorer.WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "DEV_ONLY") {
		t.Errorf("expected DEV_ONLY in unique section, got: %s", out)
	}
	if !strings.Contains(out, "dev") {
		t.Errorf("expected label 'dev' in output")
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	r := makeResult([]string{"HOST"}, map[string][]string{
		"staging": {"STAGING_KEY"},
	})
	var buf bytes.Buffer
	if err := anchorer.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["anchors"]; !ok {
		t.Error("expected 'anchors' key in JSON output")
	}
	if _, ok := out["unique"]; !ok {
		t.Error("expected 'unique' key in JSON output")
	}
}

func TestWriteJSON_EmptyAnchors(t *testing.T) {
	r := makeResult(nil, nil)
	var buf bytes.Buffer
	if err := anchorer.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), `"anchors": []`) {
		t.Errorf("expected empty anchors array, got: %s", buf.String())
	}
}
