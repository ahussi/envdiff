package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func makeEnvFile(path string, entries map[string]string) *parser.EnvFile {
	env := &parser.EnvFile{
		Path:    path,
		Entries: make(map[string]parser.Entry),
	}
	for k, v := range entries {
		env.Entries[k] = parser.Entry{Key: k, Value: v}
	}
	return env
}

func TestCompare_MissingInB(t *testing.T) {
	a := makeEnvFile(".env.a", map[string]string{"APP": "foo", "PORT": "8080"})
	b := makeEnvFile(".env.b", map[string]string{"APP": "foo"})

	result := Compare(a, b)
	if len(result.MissingInB) != 1 || result.MissingInB[0].Key != "PORT" {
		t.Errorf("expected PORT missing in B, got %+v", result.MissingInB)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := makeEnvFile(".env.a", map[string]string{"APP": "foo"})
	b := makeEnvFile(".env.b", map[string]string{"APP": "foo", "SECRET": "xyz"})

	result := Compare(a, b)
	if len(result.MissingInA) != 1 || result.MissingInA[0].Key != "SECRET" {
		t.Errorf("expected SECRET missing in A, got %+v", result.MissingInA)
	}
}

func TestCompare_Mismatch(t *testing.T) {
	a := makeEnvFile(".env.a", map[string]string{"DB": "localhost"})
	b := makeEnvFile(".env.b", map[string]string{"DB": "remotehost"})

	result := Compare(a, b)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	if result.Mismatched[0].ValueA != "localhost" || result.Mismatched[0].ValueB != "remotehost" {
		t.Errorf("unexpected mismatch values: %+v", result.Mismatched[0])
	}
}

func TestCompare_NoDifferences(t *testing.T) {
	a := makeEnvFile(".env.a", map[string]string{"APP": "foo", "PORT": "8080"})
	b := makeEnvFile(".env.b", map[string]string{"APP": "foo", "PORT": "8080"})

	result := Compare(a, b)
	if result.HasDifferences() {
		t.Error("expected no differences")
	}
}
