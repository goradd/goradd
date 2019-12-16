package strings

import "unicode"

func ExtractNumbers(in string) (out string) {
	for _,c := range in {
		if unicode.IsNumber(c) {
			out += string(c)
		}
	}
	return
}
