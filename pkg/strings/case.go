package strings

import (
	"strings"
	"unicode"
)

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

// CamelToKebab converts a camel case string to kebab case.
// If it encounters a character that is not legitimate camel case,
// it ignores it (like numbers, spaces, etc.).
// Runs of upper case letters are treated as one word.
func CamelToKebab(camelCase string) string {
	return camelToKebabOrSnake(camelCase, '-')
}

// CamelToSnake converts a camel case string to snake case.
// If it encounters a character that is not legitimate camel case,
// it ignores it (like numbers, spaces, etc.).
// Runs of upper case letters are treated as one word.
// A run of upper case, followed by lower case letters will be treated
// as if the final character in the upper case run belong with the lower case
// letters.
func CamelToSnake(camelCase string) string {
	return camelToKebabOrSnake(camelCase, '_')
}

func camelToKebabOrSnake(camelCase string, replacement rune) string {
	var kebabCase []rune
	var inUpper bool

	for i, r := range camelCase {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				if i > 0 && !inUpper {
					kebabCase = append(kebabCase, replacement)
				}
				kebabCase = append(kebabCase, unicode.ToLower(r))
				inUpper = true

			} else {
				if inUpper {
					// switching from upper to lower, if we were in an upper run
					// we need to add a hyphen in front of the last rune
					if len(kebabCase) > 1 && kebabCase[len(kebabCase)-2] != replacement {
						r2 := kebabCase[len(kebabCase)-1]
						kebabCase[len(kebabCase)-1] = replacement
						kebabCase = append(kebabCase, r2)
					}
				}
				kebabCase = append(kebabCase, r)
				inUpper = false
			}
		}
	}

	return string(kebabCase)
}
