package stringmap

import (
	"sort"
)

// These routines help range over common maps in predictable ways
// You can also copy them and substitute your own types to make custom versions of these

// SortedKeys returns the keys of any map that uses strings as keys, sorted alphabetically.
func SortedKeys[T any](m map[string]T) []string {
	keys := make([]string, len(m), len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

// Range is a convenience method to range over any map that uses strings as keys in a
// predictable order from lowest to highest. It uses
// a similar Range type function to the sync.StdMap.Range function.
func Range[T any](m map[string]T, f func(key string, val T) bool) {
	keys := SortedKeys(m)

	for _, k := range keys {
		result := f(k, m[k])
		if !result {
			break
		}
	}
}
