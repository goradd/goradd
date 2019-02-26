package control

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextbox(t *testing.T) {
	p := NewMockForm()

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
	p := NewMockForm()

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

