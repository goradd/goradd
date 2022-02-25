package iface_test

import (
	"fmt"
	"testing"

	"github.com/goradd/goradd/pkg/iface"
	"github.com/stretchr/testify/assert"
)

func ExampleIsNil() {
	fmt.Println(iface.IsNil(nil))
	// Output: true
}

type testObj struct {
	a int
}

func TestIsNil(t *testing.T) {
	var b *testObj
	var c []string
	var d map[string]interface{}
	e := map[string]interface{}{}
	var i interface{}
	var j chan bool
	var f func()

	assert.True(t, iface.IsNil(b))
	assert.False(t, iface.IsNil("s"))
	assert.True(t, iface.IsNil(b))
	assert.True(t, iface.IsNil(c))
	assert.True(t, iface.IsNil(d))
	assert.False(t, iface.IsNil(e))
	assert.True(t, iface.IsNil(i))
	assert.True(t, iface.IsNil(j))
	assert.True(t, iface.IsNil(f))
}

func TestIf(t *testing.T) {
	assert.EqualValues(t, "a", iface.If(true, "a", "b"))
	assert.EqualValues(t, "b", iface.If(false, "a", "b"))
}
