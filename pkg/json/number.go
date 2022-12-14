package json

import (
	"encoding/json"
	"strconv"
)

// NumberInt is a helper function to convert an expected integer that is returned from a json Unmarshal as a Number,
// into an actual integer without returning any errors. If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting an integer. Can convert strings too.
func NumberInt(i interface{}) int {
	switch n := i.(type) {
	case json.Number:
		v, _ := n.Int64()
		return int(v)
	case string:
		v, _ := strconv.Atoi(n)
		return v
	}
	return 0
}

// NumberFloat is a helper function to convert an expected float that is returned from a json Unmarshal as a Number,
// into an actual float64 without returning any errors. If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting a float. Can convert strings too.
func NumberFloat(i interface{}) float64 {
	switch n := i.(type) {
	case json.Number:
		v, _ := n.Float64()
		return v
	case string:
		v, _ := strconv.ParseFloat(n, 64)
		return v
	}
	return 0
}

// NumberString is a helper function to convert a value that might get cast as a Json Number into a string.
// If there is an error, it just returns 0. Use this when you absolutely
// know you are expecting a string.
func NumberString(i interface{}) string {
	switch n := i.(type) {
	case json.Number:
		v := n.String()
		return v
	case string:
		return n
	}
	panic("Unknown type for NumberString")
}
