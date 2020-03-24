package stringmap

import (
	"reflect"
	"sort"
)

// These routines help range over common maps in predictable ways
// You can also copy them and substitute your own types to make custom versions of these

// SortedKeys returns the keys of any map that uses strings as keys, sorted alphabetically.
// Note that even though we are using reflection here, this process is only slightly slower compared to not
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

// Range is a convenience method to range over any map that uses strings as keys in a
// predictable order from lowest to highest. It uses
// a similar Range type function to the sync.Map.Range function.
func Range(m interface{}, f func (key string, val interface {}) bool) {
	v := reflect.ValueOf(m)
	keys := v.MapKeys()

	sort.Slice(keys, func(a,b int) bool {
		return keys[a].String() < keys[b].String()
	})

	for _,kv := range keys {
		vv := v.MapIndex(kv)
		result := f(kv.String(), vv.Interface())
		if !result {
			break
		}
	}
}