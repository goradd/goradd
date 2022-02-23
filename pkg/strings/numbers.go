package strings

import "unicode"

// ExtractNumbers returns a string with the digits contained in the given string.
func ExtractNumbers(in string) (out string) {
	for _, c := range in {
		if unicode.IsNumber(c) {
			out += string(c)
		}
	}
	return
}
