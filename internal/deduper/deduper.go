// Package deduper identifies duplicate values across keys within one or more
// env files, helping surface copy-paste mistakes and accidental aliasing.
package deduper

import "sort"

// Entry holds a single key=value pair with its source label.
type Entry struct {
	Label string
	Key   string
	Value string
}

// Group collects all keys that share the same non-empty value.
type Group struct {
	Value string
	Keys  []Entry
}

// Result is the output of Analyse.
type Result struct {
	Groups []Group // only groups with 2+ keys are included
}

// Analyse scans the provided named env maps and returns any groups of keys
// that share an identical, non-empty value.
func Analyse(envs map[string]map[string]string) Result {
	type key struct{ label, envKey string }
	index := map[string][]Entry{}

	for label, entries := range envs {
		for k, v := range entries {
			if v == "" {
				continue
			}
			index[v] = append(index[v], Entry{Label: label, Key: k, Value: v})
		}
	}

	var groups []Group
	for val, entries := range index {
		if len(entries) < 2 {
			continue
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Label != entries[j].Label {
				return entries[i].Label < entries[j].Label
			}
			return entries[i].Key < entries[j].Key
		})
		groups = append(groups, Group{Value: val, Keys: entries})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Value < groups[j].Value
	})

	return Result{Groups: groups}
}
