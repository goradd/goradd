package strings

import "strconv"

// Convenience scripts for number conversion

// AtoUint converts a string to uint.
func AtoUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 0)
	return uint(v)
}

// AtoUint64 converts a string to uint64.
func AtoUint64(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

// AtoInt64 converts a string to Int64.
func AtoInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
