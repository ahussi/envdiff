package anchorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/anchorer"
	"github.com/user/envdiff/internal/parser"
)

func makeEnvs(pairs map[string]map[string]string) map[string]parser.EnvFile {
	out := make(map[string]parser.EnvFile, len(pairs))
	for label, kv := range pairs {
		env := make(parser.EnvFile)
		for k, v := range kv {
			env[k] = v
		}
		out[label] = env
	}
	return out
}

func TestAnalyse_AllShared(t *testing.T) {
	files := makeEnvs(map[string]map[string]string{
		"a": {"HOST": "localhost", "PORT": "8080"},
		"b": {"HOST": "prod.example.com", "PORT": "443"},
	})
	res := anchorer.Analyse(files)
	if len(res.Anchors) != 2 {
		t.Fatalf("expected 2 anchors, got %d", len(res.Anchors))
	}
}

func TestAnalyse_NoShared(t *testing.T) {
	files := makeEnvs(map[string]map[string]string{
		"a": {"ONLY_A": "1"},
		"b": {"ONLY_B": "2"},
	})
	res := anchorer.Analyse(files)
	if len(res.Anchors) != 0 {
		t.Fatalf("expected 0 anchors, got %d", len(res.Anchors))
	}
	if len(res.Unique["a"]) != 1 || res.Unique["a"][0] != "ONLY_A" {
		t.Errorf("expected ONLY_A unique to a")
	}
}

func TestAnalyse_PartialOverlap(t *testing.T) {
	files := makeEnvs(map[string]map[string]string{
		"a": {"SHARED": "x", "ONLY_A": "1"},
		"b": {"SHARED": "y", "ONLY_B": "2"},
	})
	res := anchorer.Analyse(files)
	if len(res.Anchors) != 1 || res.Anchors[0] != "SHARED" {
		t.Errorf("expected SHARED as sole anchor, got %v", res.Anchors)
	}
	if len(res.Unique["a"]) != 1 || res.Unique["a"][0] != "ONLY_A" {
		t.Errorf("unexpected unique keys for a: %v", res.Unique["a"])
	}
}

func TestAnalyse_EmptyInput(t *testing.T) {
	res := anchorer.Analyse(map[string]parser.EnvFile{})
	if len(res.Anchors) != 0 {
		t.Errorf("expected no anchors for empty input")
	}
}

func TestAnalyse_SingleFile(t *testing.T) {
	files := makeEnvs(map[string]map[string]string{
		"only": {"A": "1", "B": "2"},
	})
	res := anchorer.Analyse(files)
	// With a single file every key is an anchor.
	if len(res.Anchors) != 2 {
		t.Errorf("expected 2 anchors for single file, got %d", len(res.Anchors))
	}
	// No key is unique to a single file when there is only one file
	// (unique means count==1 AND total==1, so all qualify — verify length).
	if len(res.Unique["only"]) != 0 {
		t.Errorf("expected no unique keys when only one file present")
	}
}
