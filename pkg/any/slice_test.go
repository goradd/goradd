package any

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterfaceSlice(t *testing.T) {
	l1 := []string{"a", "b"}
	i1 := InterfaceSlice(l1)
	assert.Len(t, i1, 2)

	assert.Nil(t, InterfaceSlice(nil))

	var l2 []string
	assert.Nil(t, InterfaceSlice(l2))
}

func TestIsSlice(t *testing.T) {
	assert.True(t, IsSlice([]string{}))
	assert.False(t, IsSlice(5))
}
