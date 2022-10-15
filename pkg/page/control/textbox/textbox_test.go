package textbox

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/page"
	"github.com/stretchr/testify/assert"
)

func TestTextbox(t *testing.T) {
	p := page.NewMockForm()

	d := NewTextbox(p, "")
	d.SetText("abc")
	assert.Equal(t, d.Text(), "abc")
	assert.Equal(t, d.Value(), "abc")

	d.SetValue("defg")
	assert.Equal(t, d.Text(), "defg")
	assert.Equal(t, d.Value(), "defg")

	valid := d.MockFormValue("asdf")
	assert.Equal(t, "asdf", d.Text())
	assert.True(t, valid)
	assert.True(t, d.ValidationMessage() == "")
}

func TestTextboxValidation(t *testing.T) {
	p := page.NewMockForm()

	d := NewTextbox(p, "")
	d.SetMinLength(2)
	d.SetMaxLength(5)

	valid := d.MockFormValue("a")
	assert.Equal(t, "a", d.Text())
	assert.False(t, valid)
	assert.False(t, d.ValidationMessage() == "")

	valid = d.MockFormValue("abcdef")
	assert.Equal(t, "abcdef", d.Text())
	assert.False(t, valid)
	assert.False(t, d.ValidationMessage() == "")

	valid = d.MockFormValue("abc")
	assert.Equal(t, "abc", d.Text())
	assert.True(t, valid)
	assert.True(t, d.ValidationMessage() == "")
}

func TestExportCreatorTextbox(t *testing.T) {
	c := TextboxCreator{
		ID:          "id",
		Placeholder: "placeholder",
	}
	s := generator.ExportCreator(c)

	assert.Equal(t, `textbox.TextboxCreator{
	ID:"id",
	Placeholder:"placeholder",
}`, s)
}

func TestTextbox_Serialize(t *testing.T) {
	p := page.NewMockForm()

	c := NewTextbox(p, "")
	c.SetMinLength(2)
	c.SetMaxLength(5)

	c.MockFormValue("a")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	c.Serialize(enc)

	c2 := Textbox{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	assert.Equal(t, "a", c2.Text())
	assert.Equal(t, 2, c2.MinLength())
	assert.Equal(t, 5, c2.MaxLength())
}
