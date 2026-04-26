package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Line    int
	Comment string
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries map[string]Entry
}

// Parse reads a .env file and returns an EnvFile or an error.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:    path,
		Entries: make(map[string]Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	var pendingComment string

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			pendingComment = ""
			continue
		}

		if strings.HasPrefix(line, "#") {
			pendingComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			pendingComment = ""
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = stripQuotes(value)

		env.Entries[key] = Entry{
			Key:     key,
			Value:   value,
			Line:    lineNum,
			Comment: pendingComment,
		}
		pendingComment = ""
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", path, err)
	}

	return env, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
