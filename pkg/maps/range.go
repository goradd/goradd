package maps

import (
	"reflect"
	"sort"
)

// These routines help range over common maps in predictable ways
// You can also copy them and substitute your own types to make custom versions of these

// SortedKeys returns the keys of any map that uses strings as keys, sorted alphabetically.
// Note that even though we are using reflection here, this process is very fast compared to not
// using reflection, so feel free to use it in all situations.
func SortedKeys(i interface{}) []string {
	vMap := reflect.ValueOf(i)
	vKeys := vMap.MapKeys()
	keys := make([]string, len(vKeys), len(vKeys))
	idx := 0
	for _,vKey := range vKeys {
		keys[idx] = vKey.String()
		idx++
	}
	sort.Strings(keys)
	return keys
}
