// Package templater generates a .env.template file from one or more parsed
// env files, replacing values with empty strings or placeholder hints.
package templater

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/linter"
)

// Options controls how the template is rendered.
type Options struct {
	// Placeholders replaces values with a descriptive hint instead of empty string.
	Placeholders bool
	// IncludeComments adds a comment above each key indicating its source file.
	IncludeComments bool
}

// Entry represents a single line in the rendered template.
type Entry struct {
	Key    string
	Source string
}

// Generate builds a template from the union of keys across all provided env
// maps. keys is a map of label -> key/value pairs.
func Generate(envs map[string]map[string]string, opts Options) []Entry {
	seen := make(map[string]string) // key -> first source label
	for label, env := range envs {
		for k := range env {
			if _, exists := seen[k]; !exists {
				seen[k] = label
			}
		}
	}

	entries := make([]Entry, 0, len(seen))
	for k, src := range seen {
		entries = append(entries, Entry{Key: k, Source: src})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// Write renders the template to w.
func Write(w io.Writer, entries []Entry, opts Options) error {
	for _, e := range entries {
		if opts.IncludeComments {
			if _, err := fmt.Fprintf(w, "# source: %s\n", e.Source); err != nil {
				return err
			}
		}
		value := ""
		if opts.Placeholders {
			value = placeholderFor(e.Key)
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, value); err != nil {
			return err
		}
	}
	return nil
}

// placeholderFor returns a descriptive placeholder based on the key name.
func placeholderFor(key string) string {
	lower := strings.ToLower(key)
	switch {
	case linter.IsSensitiveKey(lower):
		return "<secret>"
	case strings.Contains(lower, "url") || strings.Contains(lower, "host"):
		return "<url>"
	case strings.Contains(lower, "port"):
		return "<port>"
	case strings.Contains(lower, "path") || strings.Contains(lower, "dir"):
		return "<path>"
	default:
		return "<value>"
	}
}
