package pinner

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeResult() Result {
	return Result{
		Pinned: []PinnedKey{
			{Key: "APP_NAME", Value: "myapp", Envs: []string{"dev", "prod"}},
		},
		Floating: []string{"DB_HOST", "LOG_LEVEL"},
	}
}

func TestWriteText_ContainsPinnedKey(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, makeResult())
	out := buf.String()
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME in text output")
	}
	if !strings.Contains(out, "myapp") {
		t.Error("expected value 'myapp' in text output")
	}
}

func TestWriteText_ContainsFloatingKeys(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, makeResult())
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in floating section")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in floating section")
	}
}

func TestWriteText_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, Result{})
	if !strings.Contains(buf.String(), "No keys found") {
		t.Error("expected 'No keys found' message for empty result")
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, makeResult()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out struct {
		Pinned []struct {
			Key   string   `json:"key"`
			Value string   `json:"value"`
			Envs  []string `json:"envs"`
		} `json:"pinned"`
		Floating []string `json:"floating"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Pinned) != 1 || out.Pinned[0].Key != "APP_NAME" {
		t.Errorf("unexpected pinned in JSON: %+v", out.Pinned)
	}
	if len(out.Floating) != 2 {
		t.Errorf("expected 2 floating in JSON, got %d", len(out.Floating))
	}
}

func TestWriteJSON_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, Result{}); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "pinned") {
		t.Error("expected 'pinned' key in JSON output")
	}
}
