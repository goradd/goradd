package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormBase_Init(t *testing.T) {

	f := &MockForm{}
	f.Self = f
	assert.Panics(t, func() {
		f.FormBase.Init(nil, "")
	})

	f = NewMockForm()
	assert.Equal(t, "MockFormID", f.ID())
}
