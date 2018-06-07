package util

import (
	"math/rand"
	"strings"
	"time"
)

const htmlValueBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789?=+-_(){}[]|#@!^"

// Generates a pseudo random string of the given length
// Characters are drawn from legal HTML values that do not need escaping
// The distribution is not perfect, so its not good for crypto, but works for general purposes
func RandomHtmlValueString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = htmlValueBytes[rand.Int63()%int64(len(htmlValueBytes))]
	}
	return string(b)
}

func EndsWith(s string, ending string) bool {
	i := strings.LastIndex(s, ending)
	return i != -1 && i == len(s)-len(ending)
}

func LcFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
