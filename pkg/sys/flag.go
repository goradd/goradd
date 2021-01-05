package sys

import (
	"github.com/goradd/goradd/pkg/strings"
	"os"
)

// GetFlagString returns the value of the given command line flag if it exists.
// This is sometimes needed to get to command line flags from init() functions.
// The value is a string, and can be in -key=value or -key value form.
func GetFlagString(key string) (string, bool) {
	for i,v := range os.Args[1:] {
		if v == key {
			if len(os.Args) > i + 2 {
				return os.Args[i + 2], true
			} else {
				return "", true
			}
		} else if strings.StartsWith(v, key + "=") {
			return v[len(key)+1:], true
		}
	}
	return "", false
}
