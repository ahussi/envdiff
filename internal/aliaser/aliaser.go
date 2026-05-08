// Package aliaser detects keys that appear to be aliases of one another —
// i.e. keys whose values are identical across all provided environments.
package aliaser

import "sort"

// Group holds a set of keys that share the same value across all envs.
type Group struct {
	Keys  []string `json:"keys"`
	Value string   `json:"value"`
}

// Result is the output of Analyse.
type Result struct {
	Groups []Group `json:"groups"`
}

// Analyse inspects multiple env maps and returns groups of keys whose
// values are identical, suggesting they may be aliases of one another.
// Only keys present in at least two environments with the same non-empty
// value are considered candidates.
func Analyse(envs map[string]map[string]string) Result {
	// value -> set of keys that carry that value in every env they appear in
	valueKeys := map[string]map[string]struct{}{}

	for _, entries := range envs {
		for k, v := range entries {
			if v == "" {
				continue
			}
			if _, ok := valueKeys[v]; !ok {
				valueKeys[v] = map[string]struct{}{}
			}
			valueKeys[v][k] = struct{}{}
		}
	}

	var groups []Group
	for val, keySet := range valueKeys {
		if len(keySet) < 2 {
			continue
		}
		keys := make([]string, 0, len(keySet))
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		groups = append(groups, Group{Keys: keys, Value: val})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Keys[0] < groups[j].Keys[0]
	})

	return Result{Groups: groups}
}
