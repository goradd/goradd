package iface_test

import (
	"fmt"
	"github.com/goradd/goradd/pkg/iface"
	"github.com/stretchr/testify/assert"
	"testing"
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

	assert.True(t, iface.IsNil(b))
}