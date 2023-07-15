package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormBase_Init(t *testing.T) {

	f := new(MockForm)
	assert.Panics(t, func() {
		f.Init(nil, "")
	})

	f = NewMockForm()
	assert.Equal(t, "MockFormID", f.ID())
}
