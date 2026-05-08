package aliaser_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/aliaser"
)

func makeEnvs(data map[string]map[string]string) map[string]map[string]string {
	return data
}

func TestAnalyse_NoAliases(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"FOO": "bar", "BAZ": "qux"},
	})
	r := aliaser.Analyse(envs)
	if len(r.Groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(r.Groups))
	}
}

func TestAnalyse_SingleAliasGroup(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"DB_PASS": "secret", "DATABASE_PASSWORD": "secret", "OTHER": "val"},
	})
	r := aliaser.Analyse(envs)
	if len(r.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(r.Groups))
	}
	if len(r.Groups[0].Keys) != 2 {
		t.Fatalf("expected 2 keys in group, got %d", len(r.Groups[0].Keys))
	}
	if r.Groups[0].Value != "secret" {
		t.Errorf("unexpected value %q", r.Groups[0].Value)
	}
}

func TestAnalyse_SkipsEmptyValues(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"X": "", "Y": ""},
	})
	r := aliaser.Analyse(envs)
	if len(r.Groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(r.Groups))
	}
}

func TestAnalyse_MultipleEnvs(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"prod": {"API_KEY": "abc123", "TOKEN": "abc123"},
		"dev":  {"API_KEY": "abc123", "TOKEN": "abc123"},
	})
	r := aliaser.Analyse(envs)
	if len(r.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(r.Groups))
	}
}

func TestWriteText_NoAliases(t *testing.T) {
	r := aliaser.Result{}
	var buf bytes.Buffer
	if err := aliaser.WriteText(&buf, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No alias") {
		t.Errorf("expected no-alias message, got %q", buf.String())
	}
}

func TestWriteText_WithAliases(t *testing.T) {
	r := aliaser.Result{
		Groups: []aliaser.Group{
			{Keys: []string{"A", "B"}, Value: "shared"},
		},
	}
	var buf bytes.Buffer
	if err := aliaser.WriteText(&buf, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "[alias]") {
		t.Errorf("expected alias label in output, got %q", buf.String())
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	r := aliaser.Result{
		Groups: []aliaser.Group{
			{Keys: []string{"X", "Y"}, Value: "val"},
		},
	}
	var buf bytes.Buffer
	if err := aliaser.WriteJSON(&buf, r); err != nil {
		t.Fatal(err)
	}
	var out aliaser.Result
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(out.Groups))
	}
}
