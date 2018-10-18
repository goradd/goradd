package util

import (
	"math/rand"
	"strings"
	"time"
)

const htmlValueBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789-_()!"

// RandomHtmlValueString generates a pseudo random string of the given length
// Characters are drawn from legal HTML values that do not need encoding.
// The distribution is not perfect, so its not good for crypto, but works for general purposes
// This also works for GET variables
func RandomHtmlValueString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = htmlValueBytes[rand.Int63()%int64(len(htmlValueBytes))]
	}
	return string(b)
}

// StartsWith returns true if the string begins with the beginning string.
func StartsWith(s string, beginning string) bool {
	if len(beginning) > len(s) {
		return false
	}
	return s[0:len(beginning)] == beginning
}

// EndsWith returns true if the string ends with the ending string.
func EndsWith(s string, ending string) bool {
	i := strings.LastIndex(s, ending)
	return i != -1 && i == len(s)-len(ending)
}

// LcFirst makes sure the first character in the string is lower case.
func LcFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
