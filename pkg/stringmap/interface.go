// package stringmap contains utilities to handle common map operations
package stringmap

import (
	"encoding/json"
	"strconv"
)

// GetString gets a string out of a map of interfaces keyed by a string. Returns the string and true
// if found, or an empty string and false if not, or the found item was not a string.
func GetString(m map[string]interface{}, key string) (string, bool) {
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

// GetJsonInt gets an integer out of a map of interfaces that were unmarshalled from json. Json unmarshalling
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

// GetJsonFloat gets a float64 out of a map of interfaces that were unmarshalled from json. Json unmarshalling
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

// GetBool gets a boolean out of a map of interfaces keyed by a string. Returns the bool and true
// if found, or false and false if not, or the found item was not a bool.
func GetBool(m map[string]interface{}, key string) (bool, bool) {
	i, ok := m[key]
	if !ok {
		return false, false
	}
	s, ok := i.(bool)
	if !ok {
		return false, false
	}
	return s, true
}



// ToStringStringMap converts a map[string]interface{} to map[string]string. Any
// items in the incoming map that are not strings are ignored.
func ToStringStringMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for k := range in {
		if s,ok := GetString(in, k); ok {
			out[k] = s
		}
	}
	return out
}