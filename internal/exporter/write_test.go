package exporter_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Kind: diff.MissingInB, ValueA: "localhost"},
		{Key: "API_KEY", Kind: diff.Mismatch, ValueA: "old", ValueB: "new"},
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Write(&buf, makeResults(), exporter.WriteOptions{
		Format: exporter.FormatText,
		FileA:  "a.env",
		FileB:  "b.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in text output, got: %s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Write(&buf, makeResults(), exporter.WriteOptions{
		Format: exporter.FormatJSON,
		FileA:  "a.env",
		FileB:  "b.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var v interface{}
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Errorf("expected valid JSON, got error: %v", err)
	}
}

func TestWrite_MarkdownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Write(&buf, makeResults(), exporter.WriteOptions{
		Format: exporter.FormatMarkdown,
		FileA:  "a.env",
		FileB:  "b.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "|") {
		t.Errorf("expected markdown table in output, got: %s", out)
	}
}

func TestWriteToFile_InfersFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	err := exporter.WriteToFile(path, makeResults(), exporter.WriteOptions{
		FileA: "a.env",
		FileB: "b.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Errorf("expected valid JSON file, got: %v", err)
	}
}

func TestWriteToFile_ExplicitFormatOverridesExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	err := exporter.WriteToFile(path, makeResults(), exporter.WriteOptions{
		Format: exporter.FormatText,
		FileA:  "a.env",
		FileB:  "b.env",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if strings.TrimSpace(string(data)) == "" {
		t.Error("expected non-empty text output")
	}
}
