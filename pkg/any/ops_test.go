package any

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIf(t *testing.T) {
	assert.Equal(t, "yes", If(true, "yes", "no"))
	assert.Equal(t, 1, If(false, 2, 1))
}

func TestIsNil(t *testing.T) {
	var m map[int]int
	assert.True(t, IsNil(nil))
	assert.True(t, IsNil(m))
	assert.False(t, IsNil(5))
}

func TestZero(t *testing.T) {
	assert.Equal(t, "", Zero[string]())
	assert.Equal(t, 0, Zero[int]())
}
