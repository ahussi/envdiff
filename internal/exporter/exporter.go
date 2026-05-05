// Package exporter writes diff results to various output file formats.
package exporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Format represents a supported export format.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "txt":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	}
	return "", fmt.Errorf("unsupported export format %q: must be one of text, json, markdown", s)
}

// InferFormat infers the export format from a file extension.
func InferFormat(path string) (Format, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "txt", "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "md", "markdown":
		return FormatMarkdown, nil
	}
	return "", fmt.Errorf("cannot infer format from extension %q", ext)
}

// OpenOutput opens a file for writing, creating it if necessary.
// The caller is responsible for closing the returned WriteCloser.
func OpenOutput(path string) (io.WriteCloser, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("exporter: open %q: %w", path, err)
	}
	return f, nil
}
