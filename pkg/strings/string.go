package strings

import (
	"strings"
)


// StartsWith returns true if the string begins with the beginning string.
func StartsWith(s string, beginning string) bool {
	return len(beginning) <= len(s) && s[:len(beginning)] == beginning
}

// EndsWith returns true if the string ends with the ending string.
func EndsWith(s string, ending string) bool {
	var sLen = len(s)
	var eLen = len(ending)

	return sLen >= eLen && s[sLen - eLen : sLen] == ending
}

// LcFirst makes sure the first character in the string is lower case.
func LcFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

