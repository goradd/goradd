package util

import (
	"fmt"
	"runtime"
)

// Various error management utilities

var MaxStackDepth = 50

// A component of a stack trace
type StackFrame struct {
	File string
	Line int
	Func string
}

// GetStackTrace returns an array of stack frames, minus "skip" frames
func GetStackTrace(skip int) (trace []StackFrame) {

	for i := 2 + skip; i < MaxStackDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		name := ""
		if f := runtime.FuncForPC(pc); f != nil {
			name = f.Name()
		}

		frame := StackFrame{file, line, name}
		trace = append(trace, frame)
	}
	return trace
}

func FormatStackTrace(trace []StackFrame) (out string) {
	for _, frame := range trace {
		out += fmt.Sprintf("%s() at %s:%d\n", frame.Func, frame.File, frame.Line)
	}
	return out
}
