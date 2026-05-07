// Package patcher generates a patch (shell export statements or .env lines)
// from a slice of diff results so that environment B can be brought in sync
// with environment A.
package patcher

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format controls the output style of the patch.
type Format string

const (
	// FormatEnv emits plain KEY=VALUE lines suitable for a .env file.
	FormatEnv Format = "env"
	// FormatExport emits `export KEY=VALUE` lines suitable for shell sourcing.
	FormatExport Format = "export"
)

// ParseFormat converts a raw string to a Format, returning an error for
// unrecognised values.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "env":
		return FormatEnv, nil
	case "export":
		return FormatExport, nil
	default:
		return "", fmt.Errorf("patcher: unknown format %q (want env|export)", s)
	}
}

// Write generates patch lines for every result whose Kind is MissingInB or
// Mismatch, writing them to w.  Keys that are only missing in A are skipped
// because the caller cannot know the correct value.
func Write(w io.Writer, results []diff.Result, format Format) error {
	for _, r := range results {
		switch r.Kind {
		case diff.MissingInB, diff.Mismatch:
			line, err := formatLine(r.Key, r.ValueA, format)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w, line); err != nil {
				return fmt.Errorf("patcher: write: %w", err)
			}
		}
	}
	return nil
}

func formatLine(key, value string, format Format) (string, error) {
	quoted := quoteValue(value)
	switch format {
	case FormatEnv:
		return fmt.Sprintf("%s=%s", key, quoted), nil
	case FormatExport:
		return fmt.Sprintf("export %s=%s", key, quoted), nil
	default:
		return "", fmt.Errorf("patcher: unsupported format %q", format)
	}
}

// quoteValue wraps value in double-quotes when it contains spaces or special
// characters; otherwise it is returned as-is.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r#") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}
