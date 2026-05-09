package deduper_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduper"
)

func makeEnvs(data map[string]map[string]string) map[string]map[string]string {
	return data
}

func TestAnalyse_NoDuplicates(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"prod": {"DB_HOST": "prod-db", "CACHE_HOST": "prod-cache"},
	})
	res := deduper.Analyse(envs)
	if len(res.Groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(res.Groups))
	}
}

func TestAnalyse_SingleDuplicateGroup(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"prod": {"DB_PASS": "secret", "API_KEY": "secret"},
	})
	res := deduper.Analyse(envs)
	if len(res.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(res.Groups))
	}
	if res.Groups[0].Value != "secret" {
		t.Errorf("unexpected value %q", res.Groups[0].Value)
	}
	if len(res.Groups[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Groups[0].Keys))
	}
}

func TestAnalyse_SkipsEmptyValues(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"dev": {"A": "", "B": ""},
	})
	res := deduper.Analyse(envs)
	if len(res.Groups) != 0 {
		t.Fatalf("empty values should not form groups, got %d", len(res.Groups))
	}
}

func TestAnalyse_AcrossMultipleEnvs(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"dev":  {"TOKEN": "abc123"},
		"prod": {"SECRET": "abc123"},
	})
	res := deduper.Analyse(envs)
	if len(res.Groups) != 1 {
		t.Fatalf("expected 1 cross-env group, got %d", len(res.Groups))
	}
	if len(res.Groups[0].Keys) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Groups[0].Keys))
	}
}

func TestAnalyse_EmptyInput(t *testing.T) {
	res := deduper.Analyse(map[string]map[string]string{})
	if len(res.Groups) != 0 {
		t.Fatalf("expected no groups for empty input")
	}
}

func TestAnalyse_SortedOutput(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"e": {"Z": "zval", "A": "aval", "M": "zval"},
	})
	res := deduper.Analyse(envs)
	if len(res.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(res.Groups))
	}
	keys := res.Groups[0].Keys
	if keys[0].Key != "M" || keys[1].Key != "Z" {
		t.Errorf("expected keys sorted: M, Z; got %s, %s", keys[0].Key, keys[1].Key)
	}
}

func TestAnalyse_MultipleGroups(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"prod": {
			"DB_PASS":  "secret",
			"API_KEY":  "secret",
			"TOKEN_A":  "tok123",
			"TOKEN_B":  "tok123",
		},
	})
	res := deduper.Analyse(envs)
	if len(res.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(res.Groups))
	}
}
