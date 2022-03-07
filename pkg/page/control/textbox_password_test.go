package control

import (
	"bytes"
	"github.com/goradd/goradd/pkg/page"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordTextbox(t *testing.T) {
	f := page.NewMockForm()
	ctx := page.NewMockContext()

	p := NewPasswordTextbox(f, "")
	assert.Equal(t, "off", p.Attributes().Get("autocomplete"))

	p.SetValue("test")
	p.SetAttribute("a", "b")

	var buf bytes.Buffer
	e := page.GobPageEncoder{}
	enc := e.NewEncoder(&buf)
	p.Serialize(enc)

	dec := e.NewDecoder(&buf)
	p2 := NewPasswordTextbox(f, "")
	p2.Deserialize(dec)
	assert.Equal(t, "", p2.value)
	assert.Equal(t, "b", p.Attributes().Get("a"))

	assert.Panics(t, func() {
		p2.SaveState(ctx, true)
	})
}

func TestPasswordCreate(t *testing.T) {
	f := page.NewMockForm()
	ctx := page.NewMockContext()
	f.AddControls(ctx,
		PasswordTextboxCreator{
		ID: "abc",
		Text: "b",
		Placeholder: "c",
		},
	)


	p := GetPasswordTextbox(f, "abc")
	assert.NotNil(t, p)
	assert.Equal(t, "abc", p.ID())
}