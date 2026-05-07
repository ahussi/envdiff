package cascader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your/envdiff/internal/cascader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestCascade_SingleLayer(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := cascader.Cascade([]cascader.Layer{{Label: "base", Path: p}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Resolved["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Resolved["FOO"])
	}
	if res.Provenance["FOO"] != "base" {
		t.Errorf("expected provenance base, got %q", res.Provenance["FOO"])
	}
}

func TestCascade_LaterLayerOverrides(t *testing.T) {
	base := writeTempEnv(t, "FOO=base_val\nSHARED=from_base\n")
	over := writeTempEnv(t, "FOO=override_val\n")
	res, err := cascader.Cascade([]cascader.Layer{
		{Label: "base", Path: base},
		{Label: "override", Path: over},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Resolved["FOO"] != "override_val" {
		t.Errorf("expected override_val, got %q", res.Resolved["FOO"])
	}
	if res.Provenance["FOO"] != "override" {
		t.Errorf("expected provenance override, got %q", res.Provenance["FOO"])
	}
	if res.Resolved["SHARED"] != "from_base" {
		t.Errorf("expected SHARED=from_base, got %q", res.Resolved["SHARED"])
	}
}

func TestCascade_SkipsEmptyPath(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	res, err := cascader.Cascade([]cascader.Layer{
		{Label: "base", Path: p},
		{Label: "missing", Path: ""},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Resolved["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", res.Resolved["KEY"])
	}
}

func TestCascade_NoLayers_ReturnsError(t *testing.T) {
	_, err := cascader.Cascade(nil)
	if err == nil {
		t.Error("expected error for empty layers, got nil")
	}
}

func TestCascade_FileNotFound_ReturnsError(t *testing.T) {
	_, err := cascader.Cascade([]cascader.Layer{
		{Label: "ghost", Path: "/nonexistent/path/.env"},
	})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
