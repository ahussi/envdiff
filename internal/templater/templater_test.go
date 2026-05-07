package templater_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/templater"
)

func makeEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		".env.production": {
			"DB_HOST":     "prod-db.example.com",
			"DB_PASSWORD": "supersecret",
			"APP_PORT":    "8080",
		},
		".env.staging": {
			"DB_HOST":  "staging-db.example.com",
			"LOG_PATH": "/var/log/app",
		},
	}
}

func TestGenerate_UnionOfKeys(t *testing.T) {
	entries := templater.Generate(makeEnvs(), templater.Options{})
	keys := make(map[string]bool)
	for _, e := range entries {
		keys[e.Key] = true
	}
	for _, expected := range []string{"DB_HOST", "DB_PASSWORD", "APP_PORT", "LOG_PATH"} {
		if !keys[expected] {
			t.Errorf("expected key %q in template entries", expected)
		}
	}
}

func TestGenerate_SortedKeys(t *testing.T) {
	entries := templater.Generate(makeEnvs(), templater.Options{})
	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted: %q before %q", entries[i-1].Key, entries[i].Key)
		}
	}
}

func TestWrite_EmptyValues(t *testing.T) {
	entries := templater.Generate(makeEnvs(), templater.Options{})
	var buf bytes.Buffer
	if err := templater.Write(&buf, entries, templater.Options{}); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	for _, line := range strings.Split(buf.String(), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || parts[1] != "" {
			t.Errorf("expected empty value, got line: %q", line)
		}
	}
}

func TestWrite_Placeholders(t *testing.T) {
	entries := templater.Generate(makeEnvs(), templater.Options{Placeholders: true})
	var buf bytes.Buffer
	if err := templater.Write(&buf, entries, templater.Options{Placeholders: true}); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "DB_PASSWORD=<secret>") {
		t.Errorf("expected DB_PASSWORD=<secret>, got:\n%s", output)
	}
	if !strings.Contains(output, "APP_PORT=<port>") {
		t.Errorf("expected APP_PORT=<port>, got:\n%s", output)
	}
}

func TestWrite_IncludeComments(t *testing.T) {
	entries := templater.Generate(makeEnvs(), templater.Options{IncludeComments: true})
	var buf bytes.Buffer
	if err := templater.Write(&buf, entries, templater.Options{IncludeComments: true}); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "# source:") {
		t.Errorf("expected source comments in output, got:\n%s", output)
	}
}
