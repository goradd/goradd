package stringmap

import (
	"fmt"
)

func ExampleSortedKeys() {
	m := map[string]int{
		"One":1,
		"Two":2,
		"Three":3,
	}

	k := SortedKeys(m)
	fmt.Print(k)
	// Output: [One Three Two]
}

