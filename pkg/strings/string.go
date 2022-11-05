package strings

import (
	"strings"
	"unicode"
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

// Indent will indent every line of the string with a tab
func Indent(s string) string {
	s = "\t" + strings.Replace(s, "\n", "\n\t", -1)
	return strings.TrimRight(s, "\t")
}

// Title is a more advanced titling operation. It will convert underscores to spaces, and add spaces to CamelCase
// words.
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
		newString += s[i : i+1]
	}
	return newString
}

// HasOnlyLetters will return false if any of the characters in the string do not pass the unicode.IsLetter test.
func HasOnlyLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// JoinContent joins strings together with the separator sep. Only strings that are not empty strings are joined.
func JoinContent(sep string, items ...string) string {
	var l []string
	for _, i := range items {
		if i != "" {
			l = append(l, i)
		}
	}
	return strings.Join(l, sep)
}

// If is like the ternary operator ?. It returns the first string on true, and the second on false.
func If(cond bool, trueVal, falseVal string) string {
	if cond {
		return trueVal
	} else {
		return falseVal
	}
}

// ContainsAnyStrings returns true if the haystack contains any of the needles
func ContainsAnyStrings(haystack string, needles ...string) bool {
	for _, h := range needles {
		if strings.Contains(haystack, h) {
			return true
		}
	}
	return false
}

// HasCharType returns true if the given string has at least one of all the
// selected char types.
func HasCharType(s string, wantUpper, wantLower, wantDigit, wantPunc, wantSymbol bool) bool {
	var hasUpper, hasLower, hasDigit, hasPunc, hasSymbol bool

	for _, c := range s {
		if !hasUpper && wantUpper && unicode.IsUpper(c) {
			hasUpper = true
		} else if !hasLower && wantLower && unicode.IsLower(c) {
			hasLower = true
		} else if !hasDigit && wantDigit && unicode.IsDigit(c) {
			hasDigit = true
		} else if !hasPunc && wantPunc && unicode.IsPunct(c) {
			hasPunc = true
		} else if !hasSymbol && wantSymbol && unicode.IsSymbol(c) {
			hasSymbol = true
		}

		if (!wantUpper || hasUpper) &&
			(!wantLower || hasLower) &&
			(!wantDigit || hasDigit) &&
			(!wantPunc || hasPunc) &&
			(!wantSymbol || hasSymbol) {
			return true
		}
	}
	return false
}
