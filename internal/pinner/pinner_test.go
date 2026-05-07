package pinner

import (
	"testing"
)

func makeEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"APP_NAME": "myapp",
			"LOG_LEVEL": "debug",
			"DB_HOST": "localhost",
		},
		"prod": {
			"APP_NAME": "myapp",
			"LOG_LEVEL": "warn",
			"DB_HOST": "db.prod.example.com",
		},
	}
}

func TestAnalyse_PinnedKey(t *testing.T) {
	r := Analyse(makeEnvs())
	if len(r.Pinned) != 1 {
		t.Fatalf("expected 1 pinned key, got %d", len(r.Pinned))
	}
	if r.Pinned[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME pinned, got %s", r.Pinned[0].Key)
	}
	if r.Pinned[0].Value != "myapp" {
		t.Errorf("unexpected pinned value: %s", r.Pinned[0].Value)
	}
}

func TestAnalyse_FloatingKeys(t *testing.T) {
	r := Analyse(makeEnvs())
	if len(r.Floating) != 2 {
		t.Fatalf("expected 2 floating keys, got %d", len(r.Floating))
	}
	expected := map[string]bool{"LOG_LEVEL": true, "DB_HOST": true}
	for _, k := range r.Floating {
		if !expected[k] {
			t.Errorf("unexpected floating key: %s", k)
		}
	}
}

func TestAnalyse_EmptyInput(t *testing.T) {
	r := Analyse(map[string]map[string]string{})
	if len(r.Pinned) != 0 || len(r.Floating) != 0 {
		t.Error("expected empty result for empty input")
	}
}

func TestAnalyse_SingleEnv(t *testing.T) {
	// With only one env, no key can be pinned across multiple envs.
	r := Analyse(map[string]map[string]string{
		"dev": {"APP_NAME": "myapp", "LOG_LEVEL": "debug"},
	})
	if len(r.Pinned) != 0 {
		t.Errorf("expected 0 pinned keys for single env, got %d", len(r.Pinned))
	}
	if len(r.Floating) != 2 {
		t.Errorf("expected 2 floating keys for single env, got %d", len(r.Floating))
	}
}

func TestAnalyse_AllPinned(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"X": "1", "Y": "2"},
		"b": {"X": "1", "Y": "2"},
		"c": {"X": "1", "Y": "2"},
	}
	r := Analyse(envs)
	if len(r.Pinned) != 2 {
		t.Errorf("expected 2 pinned keys, got %d", len(r.Pinned))
	}
	if len(r.Floating) != 0 {
		t.Errorf("expected 0 floating keys, got %d", len(r.Floating))
	}
}
