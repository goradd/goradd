package control

import (
	"github.com/goradd/goradd/pkg/page/event"
	"testing"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/stretchr/testify/assert"
)

func TestNewButton(t *testing.T) {
	form := page.NewMockForm()
	ctx := page.NewMockContext()

	b := NewButton(form, "btnId")
	assert.Equal(t, "button", b.Tag)

	b.SetLabel("mybutton")
	assert.Equal(t, "mybutton", b.Text())

	b.SetIsPrimary(true)
	assert.Equal(t, "submit", b.Attributes().Get("type"))

	b.SetIsPrimary(false)
	assert.Equal(t, "button", b.Attributes().Get("type"))

	a := b.DrawingAttributes(ctx)
	assert.True(t, a.Has("name"))
	assert.Equal(t, b.ID(), a.Get("value"))
}

func TestButtonOn(t *testing.T) {
	form := page.NewMockForm()
	//ctx := page.NewMockContext()

	b := NewButton(form, "btnId")

	assert.Nil(t, b.Event("click"))
	b.OnSubmit(action.Confirm("Help!"))
	assert.NotNil(t, b.Event("click"))
	b.Off()
	assert.Nil(t, b.Event("click"))
}

func TestButtonCreator(t *testing.T) {
	f := page.NewMockForm()
	ctx := page.NewMockContext()
	f.AddControls(ctx,
		ButtonCreator{
			ID:             "abc",
			Text:           "b",
			IsPrimary:      true,
			OnClick:        action.Blur("d"),
			ValidationType: event.ValidateChildrenOnly,
			ControlOptions: page.ControlOptions{
				Class: "c",
			},
		},
	)

	b := GetButton(f, "abc")
	assert.NotNil(t, b)
	assert.Equal(t, "abc", b.ID())
	assert.True(t, b.HasClass("c"))
	assert.Equal(t, "submit", b.Attributes().Get("type"))
	assert.NotNil(t, b.Event("click"))
}
