package log

import (
	"fmt"
	"runtime"
)

// StackTrace returns a formatted stack trace, listing files and functions on each line.
func StackTrace(startingDepth int, maxDepth int) (out string) {
	for i := 1 + startingDepth; i < maxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		name := ""
		if f := runtime.FuncForPC(pc); f != nil {
			name = f.Name()
		}

		// This format allows IDEs to turn the listing into links to the file and line.
		out += fmt.Sprintf("%s:%d -- %s\n", file, line, name)
	}
	return
}
