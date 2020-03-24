package strings

import "strings"

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

func SnakeToKebab(s string) string {
	return strings.Replace(s, "_", "-", -1)
}


