package strings

import "unicode"

// Routines in this file are aids to validation checking

func IsASCII(s string) bool {
	// idea adapted from here:
	// https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/unicode/utf8/utf8.go;l=528
	for len(s) > 0 {
		if len(s) >= 8 {
			first32 := uint32(s[0]) | uint32(s[1])<<8 | uint32(s[2])<<16 | uint32(s[3])<<24
			second32 := uint32(s[4]) | uint32(s[5])<<8 | uint32(s[6])<<16 | uint32(s[7])<<24
			if (first32|second32)&0x80808080 != 0 {
				return false
			}
			s = s[8:]
			continue
		}
		if s[0] > unicode.MaxASCII {
			return false
		}
		s = s[1:]
	}
	return true
}
