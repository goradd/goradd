package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControlBase_Init(t *testing.T) {
	t.Run("normal create", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Init(b, f, "testid")

		assert.True(t, b.NeedsRefresh())
		assert.Equal(t, "testid", b.ID())
	})

	t.Run("id with underscore", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		assert.Panics(t, func() {
			b.Init(b, f, "test_id")
		})
	})

	t.Run("same id", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Init(b, f, "testid")
		b2 := new(ControlBase)
		assert.Panics(t, func() {
			b2.Init(b, f, "testid")
		})
	})

	t.Run("no id", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Init(b, f, "")
		assert.Equal(t, "c1", b.ID())
	})
}
