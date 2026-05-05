package exporter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/exporter"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  exporter.Format
	}{
		{"text", exporter.FormatText},
		{"txt", exporter.FormatText},
		{"TEXT", exporter.FormatText},
		{"json", exporter.FormatJSON},
		{"JSON", exporter.FormatJSON},
		{"markdown", exporter.FormatMarkdown},
		{"md", exporter.FormatMarkdown},
		{"Markdown", exporter.FormatMarkdown},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := exporter.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := exporter.ParseFormat("csv")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestInferFormat(t *testing.T) {
	cases := []struct {
		path string
		want exporter.Format
	}{
		{"output.txt", exporter.FormatText},
		{"output.json", exporter.FormatJSON},
		{"output.md", exporter.FormatMarkdown},
		{"output.markdown", exporter.FormatMarkdown},
		{"output", exporter.FormatText},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got, err := exporter.InferFormat(tc.path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestInferFormat_Unknown(t *testing.T) {
	_, err := exporter.InferFormat("output.csv")
	if err == nil {
		t.Fatal("expected error for unknown extension")
	}
}

func TestOpenOutput_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	wc, err := exporter.OpenOutput(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestOpenOutput_InvalidPath(t *testing.T) {
	_, err := exporter.OpenOutput("/nonexistent/dir/out.txt")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
