package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackTrace(t *testing.T) {
	s := StackTrace(0, 5)
	assert.Contains(t, s, "stack_test.go")
}
