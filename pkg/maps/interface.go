// package maps contains utilities to handle common map operations
package maps

import (
	"encoding/json"
	"strconv"
)

// IMapGetString gets a string out of a map of interfaces keyed by a string. Returns the string and true
// if found, or an empty string and false if not, or the found item was not a string.
func GetIString(m map[string]interface{}, key string) (string, bool) {
	i, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := i.(string)
	if !ok {
		return "", false
	}
	return s, true
}

// IMapGetJsonInt gets an integer out of a map of interfaces that were unmarshalled from json. Json unmarshalling
// creates json numbers, and these have to be coerced into integers or floats.
// Returns an int and true if found, and a zero and false if not, or the found item was not a coercible number.
func GetJsonInt(m map[string]interface{}, key string) (int, bool) {
	i, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := i.(type) {
	case json.Number:
		v, _ := n.Int64()
		return int(v), true
	case string:
		v, _ := strconv.Atoi(n)
		return v, true
	}
	return 0, false
}

// IMapGetJsonFloat gets a float64 out of a map of interfaces that were unmarshalled from json. Json unmarshalling
// creates json numbers, and these have to be coerced into integers or floats.
// Returns an int and true if found, and a zero and false if not, or the found item was not a coercible number.
func GetJsonFloat(m map[string]interface{}, key string) (float64, bool) {
	i, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := i.(type) {
	case json.Number:
		v, _ := n.Float64()
		return v, true
	case string:
		v, _ := strconv.ParseFloat(n, 64)
		return v, true
	}
	return 0, false
}
