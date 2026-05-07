package cascader_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your/envdiff/internal/cascader"
)

func makeResult() *cascader.Result {
	return &cascader.Result{
		Resolved:   map[string]string{"APP_ENV": "production", "DB_URL": "postgres://localhost/db"},
		Provenance: map[string]string{"APP_ENV": "base", "DB_URL": "override"},
		Layers: []cascader.Layer{
			{Label: "base", Path: ".env"},
			{Label: "override", Path: ".env.prod"},
		},
	}
}

func TestWriteText_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	cascader.WriteText(&buf, makeResult())
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
	if !strings.Contains(out, "DB_URL") {
		t.Error("expected DB_URL in output")
	}
}

func TestWriteText_ContainsProvenance(t *testing.T) {
	var buf bytes.Buffer
	cascader.WriteText(&buf, makeResult())
	out := buf.String()
	if !strings.Contains(out, "override") {
		t.Error("expected provenance label 'override' in output")
	}
}

func TestWriteText_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	cascader.WriteText(&buf, &cascader.Result{
		Resolved: map[string]string{},
	})
	if !strings.Contains(buf.String(), "No keys resolved") {
		t.Error("expected 'No keys resolved' message")
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := cascader.WriteJSON(&buf, makeResult()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["resolved"]; !ok {
		t.Error("expected 'resolved' key in JSON")
	}
	if _, ok := out["provenance"]; !ok {
		t.Error("expected 'provenance' key in JSON")
	}
	if _, ok := out["layers"]; !ok {
		t.Error("expected 'layers' key in JSON")
	}
}
