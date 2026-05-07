// Package anchorer detects keys that are present in all provided env files,
// treating them as "anchors" — a stable baseline shared across environments.
package anchorer

import (
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the outcome of an anchor analysis.
type Result struct {
	// Anchors are keys present in every env file.
	Anchors []string
	// Unique maps each label to keys found only in that file.
	Unique map[string][]string
}

// Analyse compares multiple named env files and returns keys that are common
// to all of them (anchors) and keys exclusive to each file.
func Analyse(files map[string]parser.EnvFile) Result {
	if len(files) == 0 {
		return Result{Unique: make(map[string][]string)}
	}

	// Count how many files each key appears in.
	count := make(map[string]int)
	for _, env := range files {
		for key := range env {
			count[key]++
		}
	}

	total := len(files)
	anchorSet := make(map[string]bool)
	for key, c := range count {
		if c == total {
			anchorSet[key] = true
		}
	}

	// Build sorted anchor slice.
	anchors := make([]string, 0, len(anchorSet))
	for key := range anchorSet {
		anchors = append(anchors, key)
	}
	sort.Strings(anchors)

	// Build unique map: keys that appear in exactly one file.
	unique := make(map[string][]string)
	for label, env := range files {
		for key := range env {
			if count[key] == 1 {
				unique[label] = append(unique[label], key)
			}
		}
		sort.Strings(unique[label])
	}

	return Result{
		Anchors: anchors,
		Unique:  unique,
	}
}
