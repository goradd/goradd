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

	return sLen >= eLen && s[sLen-eLen:sLen] == ending
}

// LcFirst makes sure the first character in the string is lower case.
func LcFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// Indent will indent very line of the string with a tab
func Indent(s string) string {
	s = "\t" + strings.Replace(s, "\n", "\n\t", -1)
	return strings.TrimRight(s, "\t")
}

// Title is a more advanced titling operation, that will convert underscores to spaces, and add spaces to CamelCase
// words
func Title(s string) string {
	s = strings.TrimSpace(strings.Title(strings.Replace(s, "_", " ", -1)))
	if len(s) <= 1 {
		return s
	}

	newString := s[0:1]
	l := strings.ToLower(s)
	for i := 1; i < len(s); i++ {
		if l[i] != s[i] && s[i-1:i] != " " {
			newString += " "
		}
		newString += s[i:i+1]
	}
	return newString
}

func KebabToCamel(s string) string {
	var r string

	for _, w := range strings.Split(s, "-") {
		if w == "" {
			continue
		}

		u := strings.Title(w)
		r += u
	}

	return r
}