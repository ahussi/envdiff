package grapher

import (
	"testing"
)

func makeEnvs(data map[string]map[string]string) map[string]map[string]string {
	return data
}

func TestAnalyse_TwoIdenticalMaps(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"FOO": "1", "BAR": "2"},
		"b": {"FOO": "1", "BAR": "2"},
	})
	r := Analyse(envs)
	if len(r.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(r.Edges))
	}
	e := r.Edges[0]
	if e.SharedKeys != 2 {
		t.Errorf("shared=%d want 2", e.SharedKeys)
	}
	if e.Similarity != 1.0 {
		t.Errorf("similarity=%.2f want 1.00", e.Similarity)
	}
}

func TestAnalyse_DisjointMaps(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"FOO": "1"},
		"b": {"BAR": "2"},
	})
	r := Analyse(envs)
	e := r.Edges[0]
	if e.SharedKeys != 0 {
		t.Errorf("shared=%d want 0", e.SharedKeys)
	}
	if e.Similarity != 0.0 {
		t.Errorf("similarity=%.2f want 0.00", e.Similarity)
	}
}

func TestAnalyse_PartialOverlap(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"FOO": "1", "BAR": "2", "BAZ": "3"},
		"b": {"FOO": "1", "QUX": "4"},
	})
	r := Analyse(envs)
	e := r.Edges[0]
	if e.SharedKeys != 1 {
		t.Errorf("shared=%d want 1", e.SharedKeys)
	}
	// union = 3+2-1 = 4, sim = 1/4 = 0.25
	if e.Similarity < 0.24 || e.Similarity > 0.26 {
		t.Errorf("similarity=%.4f want ~0.25", e.Similarity)
	}
}

func TestAnalyse_ThreeNodes_EdgeCount(t *testing.T) {
	envs := makeEnvs(map[string]map[string]string{
		"a": {"X": "1"},
		"b": {"Y": "2"},
		"c": {"Z": "3"},
	})
	r := Analyse(envs)
	if len(r.Edges) != 3 {
		t.Errorf("edges=%d want 3", len(r.Edges))
	}
	if len(r.Nodes) != 3 {
		t.Errorf("nodes=%d want 3", len(r.Nodes))
	}
}

func TestAnalyse_EmptyInput(t *testing.T) {
	r := Analyse(map[string]map[string]string{})
	if len(r.Edges) != 0 {
		t.Errorf("expected no edges for empty input")
	}
}
