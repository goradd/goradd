package strings

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Routines in this file are aids to validation checking

// IsASCII returns true if the string contains only ascii characters
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

// IsUTF8 returns true if the given string only contains valid UTF-8 characters
func IsUTF8(s string) bool {
	return utf8.ValidString(s)
}

// IsUTF8Bytes returns true if the given byte array only contains valid UTF-8 characters
func IsUTF8Bytes(b []byte) bool {
	return utf8.Valid(b)
}

// IsInt returns true if the given string is an integer.
// Allows the string to start with a + or -.
func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// IsFloat returns true if the given string is a floating point number.
// Allows the string to start with a + or -.
func IsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func StripNewlines(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func StripNulls(s string) string {
	s = strings.Replace(s, "\000", "", -1)
	return s
}

func HasNull(s string) bool {
	return strings.Contains(s, "\000")
}
