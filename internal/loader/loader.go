package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a parsed environment file with its source path.
type EnvFile struct {
	Path   string
	Name   string
	Values map[string]string
}

// Load reads and parses a .env file from the given path.
func Load(path string) (*EnvFile, error) {
	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("cannot access file %s: %w", path, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file: %s", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer f.Close()

	values, err := parser.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("cannot parse file %s: %w", path, err)
	}

	return &EnvFile{
		Path:   path,
		Name:   labelFromPath(path),
		Values: values,
	}, nil
}

// LoadPair loads two env files and returns them as a pair.
func LoadPair(pathA, pathB string) (*EnvFile, *EnvFile, error) {
	a, err := Load(pathA)
	if err != nil {
		return nil, nil, fmt.Errorf("loading first file: %w", err)
	}

	b, err := Load(pathB)
	if err != nil {
		return nil, nil, fmt.Errorf("loading second file: %w", err)
	}

	return a, b, nil
}

// labelFromPath derives a short display name from a file path.
func labelFromPath(path string) string {
	base := filepath.Base(path)
	// Strip leading dot for display, e.g. ".env.production" -> "env.production"
	return strings.TrimPrefix(base, ".")
}
