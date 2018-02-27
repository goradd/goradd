package control
import (
	//"github.com/microcosm-cc/bluemonday"
	"net/mail"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
)

const (
	TEXTBOX_TYPE_DEFAULT    = "text"
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

	sanitizer Sanitizer
	ValidationFilter func(string)bool
	ValidationMessage string

	minLength int
	maxLength int
}

// Creates a new standard html text box
func NewTextBox(parent page.ControlI, id string) *TextBox {
	t := &TextBox{}
	t.Init(t, parent, id)
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
	t.SetAttribute("type", TEXTBOX_TYPE_DEFAULT)
}

func (t *TextBox) this() TextBoxI {
	return t.Self.(TextBoxI)
}

func (t *TextBox) Attributes() *html.Attributes {
	attr := html.NewAttributes()
	attr.Set("name", t.Id())	// needed for posts
	a := t.Control.Attributes()
	return a.Override(attr)
}

// Set the value of the text. Returns itself for chaining
func (t *TextBox) SetText(s string) page.ControlI {
	t.SetAttribute("value", s)
	return t.this()
}

func (t *TextBox) Text() string {
	return t.Attribute("value")
}

func (t *TextBox) SetValue(v interface{}) page.ControlI {
	return t.SetText(v.(string))
}

func (t *TextBox) Value() interface{} {
	return t.Text()
}

func (t *TextBox) SetMaxLength(len int) TextBoxI {
	t.maxLength = len
	t.SetAttribute("maxlength", len)
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

// SetType sets the type of textbox this is. Pass it a TEXTBOX_TYPE... constant
func (t *TextBox) SetType(s string) TextBoxI {
	t.Attributes().Set("type", s) // directly set our attribute value to reflect what the browser already has
	t.Refresh() // can't change this without completely redrawing the control
	return t.this()
}

func (t *TextBox) ParsePostData(c context.Context) {
	ctx := page.GetContext(c)
	if text,ok := ctx.FormValue(t.Id()); ok {
		text = t.sanitize(text)
		t.Attributes().Set("value", text) // directly set our attribute value to reflect what the browser already has
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