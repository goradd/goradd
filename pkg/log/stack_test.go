package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackTrace(t *testing.T) {
	s := StackTrace(0,5)
	assert.Contains(t, s, "stack_test.go")
}
