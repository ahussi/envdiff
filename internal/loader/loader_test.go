package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoad_BasicFile(t *testing.T) {
	path := writeTempEnv(t, ".env", "FOO=bar\nBAZ=qux\n")

	ef, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ef.Values["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", ef.Values["FOO"])
	}
	if ef.Values["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", ef.Values["BAZ"])
	}
	if ef.Name != "env" {
		t.Errorf("expected name 'env', got %q", ef.Name)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_Directory(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error when path is a directory, got nil")
	}
}

func TestLoadPair_Success(t *testing.T) {
	pathA := writeTempEnv(t, ".env", "KEY=one\n")
	pathB := writeTempEnv(t, ".env.production", "KEY=two\n")

	a, b, err := LoadPair(pathA, pathB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Values["KEY"] != "one" {
		t.Errorf("expected a.KEY=one, got %q", a.Values["KEY"])
	}
	if b.Values["KEY"] != "two" {
		t.Errorf("expected b.KEY=two, got %q", b.Values["KEY"])
	}
}

func TestLoadPair_SecondFileMissing(t *testing.T) {
	pathA := writeTempEnv(t, ".env", "KEY=one\n")
	_, _, err := LoadPair(pathA, "/no/such/file")
	if err == nil {
		t.Fatal("expected error when second file is missing, got nil")
	}
}

func TestLabelFromPath(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"/some/dir/.env", "env"},
		{"/some/dir/.env.staging", "env.staging"},
		{"/some/dir/config.env", "config.env"},
	}
	for _, tc := range cases {
		got := labelFromPath(tc.input)
		if got != tc.want {
			t.Errorf("labelFromPath(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
