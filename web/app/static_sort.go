package app

import (
	"github.com/goradd/gengen/pkg/maps"
	"sort"
)

// Sort by keys interface
type sortStringbykeylen struct {
	// This embedded interface permits Reverse to use the methods of
	// another interface implementation.
	sort.Interface
}

// A helper function to allow StringSliceMaps to be sorted by keys
// To sort the map by keys, call:
//   sort.Sort(OrderStringStringSliceMapByKeys(m))
func OrderDirectoryPaths(o *maps.StringSliceMap) sort.Interface {
	return &sortStringbykeylen{o}
}

// A helper function to allow StringSliceMaps to be sorted by keys
func (r sortStringbykeylen) Less(i, j int) bool {
	var o *maps.StringSliceMap = r.Interface.(*maps.StringSliceMap)

	// order longest to shortest
	return len(o.GetKeyAt(i)) > len(o.GetKeyAt(j))
}
