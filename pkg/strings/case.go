package strings

import "strings"

// KebabToCamel convert kebab-case words to CamelCase words.
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

// SnakeToKebab converts snake_case words to kebab-case words.
func SnakeToKebab(s string) string {
	return strings.Replace(s, "_", "-", -1)
}
