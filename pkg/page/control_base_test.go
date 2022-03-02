package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControlBase_Init(t *testing.T) {
	t.Run("normal create", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Self = b
		b.Init(f, "testid")

		assert.True(t, b.NeedsRefresh())
		assert.Equal(t, "testid", b.ID())
	})

	t.Run("id with underscore", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Self = b
		assert.Panics(t, func() {
			b.Init(f, "test_id")
		})
	})

	t.Run("same id", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Self = b
		b.Init(f, "testid")
		b2 := new(ControlBase)
		b2.Self = b2
		assert.Panics(t, func() {
			b2.Init(f, "testid")
		})
	})

	t.Run("no id", func(t *testing.T) {
		f := NewMockForm()
		b := new(ControlBase)
		b.Self = b
		b.Init(f, "")
		assert.Equal(t, "c1", b.ID())
	})
}
