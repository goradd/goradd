package control
import (
	//"github.com/microcosm-cc/bluemonday"
	"net/mail"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
	"strconv"
)

const (
	TEXTBOX_TYPE_DEFAULT   = "text"
	TEXTBOX_TYPE_PASSWORD    = "password"
	TEXTBOX_TYPE_SEARCH    = "search"
	TEXTBOX_TYPE_NUMBER = "number"
	TEXTBOX_TYPE_EMAIL = "email"
	TEXTBOX_TYPE_TEL = "tel"
	TEXTBOX_TYPE_URL = "url"
)

type Sanitizer interface {
	Sanitize(string)string
}

type TextBoxI interface {
	page.ControlI
}

type TextBox struct {
	localPage.Control

	typ string

	sanitizer Sanitizer
	ValidationFilter func(string)bool
	ValidationMessage string

	minLength int
	maxLength int

	value string
}

// Creates a new standard html text box
func NewTextBox(parent page.ControlI, id string) *TextBox {
	t := &TextBox{}
	t.Init(t, parent, id)
	t.typ = TEXTBOX_TYPE_DEFAULT
	return t
}

//TODO
func (t *TextBox) Serialize(buf []byte) {

}

func (t *TextBox) Unserialize(data interface{}) {

}

// Initializes a textbox. Normally you will not call this directly. However, sub controls should call this after
// creation to get the enclosed control initialized. Self is the newly created class. Like so:
// t := &MyTextBox{}
// t.TextBox.Init(t, parent, id)
// A parent control is isRequired. Leave id blank to have the system assign an id to the control.
func (t *TextBox) Init(self TextBoxI, parent page.ControlI, id string) {
	t.Control.Init(self, parent, id)

	t.Tag = "input"
	t.IsVoidTag = true
	t.typ = TEXTBOX_TYPE_DEFAULT
}

func (t *TextBox) this() TextBoxI {
	return t.Self.(TextBoxI)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (t *TextBox) DrawingAttributes() *html.Attributes {
	a := t.Control.DrawingAttributes()
	a.Set("name", t.Id())	// needed for posts
	a.Set("type", t.typ)
	a.Set("value", t.value)
	if t.maxLength != 0 {
		a.Set("maxlength", strconv.Itoa(t.maxLength))
	}
	if t.minLength != 0 {
		a.Set("minlength", strconv.Itoa(t.maxLength))
	}

	return a
}

// Set the value of the text. Returns itself for chaining
func (t *TextBox) SetText(s string) page.ControlI {
	t.SetAttribute("value", s)
	return t.this()
}

func (t *TextBox) Text() string {
	return t.value
}

func (t *TextBox) SetValue(v interface{}) page.ControlI {
	t.value,_ = v.(string)
	return t.this()
}

func (t *TextBox) Value() interface{} {
	return t.Text()
}

func (t *TextBox) SetMaxLength(len int) TextBoxI {
	t.maxLength = len
	return t.this()
}

func (t *TextBox) MaxLength() int {
	return t.maxLength
}

func (t *TextBox) SetMinLength(len int) TextBoxI {
	t.minLength = len
	return t.this()
}

func (t *TextBox) MinLength() int {
	return t.minLength
}

func (t *TextBox) SetPlaceholder(s string) TextBoxI {
	t.SetAttribute("placeholder", s)
	return t.this()
}

func (t *TextBox) Placeholder() string {
	return t.Attribute("placeholder")
}

// SetType sets the type of textbox this is. Pass it a TEXTBOX_TYPE... constant normally, thought you can pass
// any string and it will become the input type
func (t *TextBox) SetType(typ string) TextBoxI {
	t.typ = typ
	t.Refresh() // can't change this without completely redrawing the control
	return t.this()
}

func (t *TextBox) ParsePostData(c context.Context) {
	ctx := page.GetContext(c)
	if text,ok := ctx.FormValue(t.Id()); ok {
		t.value = t.sanitize(text)
	}
}

func (t *TextBox) SetSanitizer(s Sanitizer) {
	t.sanitizer = s
}
func (t *TextBox) sanitize(s string) string {
	if t.sanitizer == nil {
		t.sanitizer = DefaultSanitizer{}
	}
	return t.sanitizer.Sanitize(s)
}

// Sanitizers

type DefaultSanitizer struct {
}

func (d DefaultSanitizer) Sanitize(s string) string {
	/*
	s = strings.TrimSpace(s)
	p := bluemonday.StrictPolicy()
	s = p.Sanitize(s)
	return s*/
	return ""
}

// Validators
func (t *TextBox) Validate()bool {
	valid := true
	text := t.Text()
	if t.Required() && text == "" {
		valid = false
		t.SetValidationError("Value is isRequired")
	}
	if t.ValidationFilter != nil && !t.ValidationFilter(text) {
		valid = false
		t.SetValidationError(t.ValidationMessage)
	}
	return valid
}

func ValidateEmail(s string)bool {
	_, err := mail.ParseAddressList(s)
	if err != nil {
		return false
	}
	return true
}