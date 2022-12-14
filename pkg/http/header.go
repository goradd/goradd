package http

import "strings"

// ParseValueAndParams returns the value and param map for Content-Type and Content-Disposition header values
func ParseValueAndParams(in string) (value string, params map[string]string) {
	parts := strings.Split(in, ";")
	if len(parts) > 0 {
		value = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			for _, p := range parts[1:] {
				p = strings.TrimSpace(p)
				offset := strings.IndexRune(p, '=')
				if offset >= 0 {
					if params == nil {
						params = make(map[string]string)
					}
					params[p[:offset]] = p[offset+1:]
				}
			}
		}
	}
	return
}
