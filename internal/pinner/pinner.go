// Package pinner identifies keys whose values are identical across all
// provided environments ("pinned" values) and keys that vary ("floating").
package pinner

import "sort"

// Result holds the output of a pin analysis.
type Result struct {
	// Pinned contains keys whose value is the same in every environment.
	Pinned []PinnedKey
	// Floating contains keys whose value differs across at least two environments.
	Floating []string
}

// PinnedKey is a key with its shared value and the environments it was found in.
type PinnedKey struct {
	Key    string
	Value  string
	Envs   []string
}

// Analyse inspects a map of environment-label → key/value pairs and returns
// which keys are pinned (identical everywhere they appear) and which float.
func Analyse(envs map[string]map[string]string) Result {
	if len(envs) == 0 {
		return Result{}
	}

	// Collect all keys across all envs.
	keyEnvs := map[string]map[string]string{} // key → { envLabel → value }
	for label, entries := range envs {
		for k, v := range entries {
			if keyEnvs[k] == nil {
				keyEnvs[k] = map[string]string{}
			}
			keyEnvs[k][label] = v
		}
	}

	var pinned []PinnedKey
	var floating []string

	for key, envValues := range keyEnvs {
		if len(envValues) < 2 {
			// Only appears in one env — treat as floating (not universally pinned).
			for label := range envValues {
				_ = label
			}
			floating = append(floating, key)
			continue
		}

		var firstVal string
		var labels []string
		uniform := true
		for label, val := range envValues {
			labels = append(labels, label)
			if firstVal == "" && len(labels) == 1 {
				firstVal = val
			}
			if val != firstVal {
				uniform = false
			}
		}

		if uniform {
			sort.Strings(labels)
			pinned = append(pinned, PinnedKey{Key: key, Value: firstVal, Envs: labels})
		} else {
			floating = append(floating, key)
		}
	}

	sort.Slice(pinned, func(i, j int) bool { return pinned[i].Key < pinned[j].Key })
	sort.Strings(floating)

	return Result{Pinned: pinned, Floating: floating}
}
