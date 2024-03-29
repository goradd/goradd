package page_test

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/textbox"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
func TestEmptyFormEncoding(t *testing.T) {
	var form = page.FormBase{}
	var b bytes.Buffer

	gob.Register(&form) // register form here, since we normally would never register the FormBase
	form.Init(nil, &form, "", "TestForm")


	pe := page.GobPageEncoder{}
	e := pe.NewEncoder(&b)

	assert.NoError(t, e.Serialize(form.Page()))

	d := pe.NewDecoder(bytes.NewBuffer(b.Bytes()))

	var p2 page.Page

	assert.NoError(t, d.Deserialize(&p2))

	assert.Equal (t, "TestForm", p2.Form().ID(), "Form id not restored")
}*/

type BasicForm struct {
	page.MockForm

	S string
}

func CreateBasicForm(ctx context.Context) page.FormI {
	f := new(BasicForm)
	f.FormBase.Init(f, ctx, "BasicForm")
	f.CreateControls(ctx)
	return f
}

func (f *BasicForm) CreateControls(ctx context.Context) {
	textbox.NewTextbox(f, "txt1").SetValue("Hi")
	f.S = "test"
}

func TestBasicFormEncoding(t *testing.T) {
	var form = CreateBasicForm(nil)
	var b bytes.Buffer

	page.RegisterControl(&BasicForm{})

	pe := page.GobPageEncoder{}
	e := pe.NewEncoder(&b)

	assert.NoError(t, e.Encode(form.Page()))

	d := pe.NewDecoder(bytes.NewBuffer(b.Bytes()))

	var p2 page.Page

	assert.NoError(t, d.Decode(&p2))

	f2 := p2.Form().(*BasicForm)

	assert.Equal(t, "BasicForm", f2.ID(), "Form id not restored")

	assert.Equal(t, "Hi", textbox.GetTextbox(f2, "txt1").Text(), "Textbox content not restored")
	assert.Equal(t, "test", f2.S)
}
