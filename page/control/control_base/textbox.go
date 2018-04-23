package control_base
import (
	//"github.com/microcosm-cc/bluemonday"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
	"strconv"
	"github.com/spekary/goradd/util/types"
	"fmt"
	"bytes"
	"context"
	html2 "html"
)

type Sanitizer interface {
	Sanitize(string)string
}

// A TextboxValidator can be added to a textbox to validate its input on the server side. A textbox can have more than one validator.
// A number of built-in validators are provided.
type Validater interface {
	// Validate evaluates the input, and returns an empty string if the input is valid, and an error string to display
	// to the user if the input does not pass the validator.
	Validate(page.Translater, string) (string)
}

type TextboxI interface {
	page.ControlI
}

type Textbox struct {
	localPage.Control

	typ string

	sanitizer Sanitizer
	validators []Validater

	minLength int
	maxLength int

	value string

	columns int
	rows int

	readonly bool
}




// Initializes a textbox. Normally you will not call this directly. However, sub controls should call this after
// creation to get the enclosed control initialized. Self is the newly created class. Like so:
// t := &MyTextBox{}
// t.Textbox.Init(t, parent, id)
// A parent control is isRequired. Leave id blank to have the system assign an id to the control.
func (t *Textbox) Init(self TextboxI, parent page.ControlI) {
	t.Control.Init(self, parent)

	t.Tag = "input"
	t.IsVoidTag = true
	t.typ = "text" // default
	t.SetHasFor(true)
}

func (t *Textbox) this() TextboxI {
	return t.Self.(TextboxI)
}

// ValidateWith adds a TextboxValidator to the validator list
func (t *Textbox) ValidateWith (v Validater) {
	t.validators = append(t.validators, v)
}

func (t *Textbox) ResetValidators () {
	t.validators = nil
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (t *Textbox) DrawingAttributes() *html.Attributes {
	a := t.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "textbox")
	a.Set("name", t.Id())	// needed for posts
	if t.Required() {
		a.Set("required", "")
	}
	if t.maxLength != 0 {
		a.Set("maxlength", strconv.Itoa(t.maxLength))
	}
	if t.rows == 0 {	// single-line textbox
		a.Set("type", t.typ)
		a.Set("value", t.value)
		if t.columns != 0 {
			a.Set("size", strconv.Itoa(t.columns))
		}
	} else {
		a.Set("rows", strconv.Itoa(t.rows))
		if t.columns != 0 {
			a.Set("cols", strconv.Itoa(t.columns))
		}
	}
	return a
}

// DrawInnerHtml is an internal function that renders the inner html of a tag. In this case, it is rendering the inner
// text of a textarea
func (t *Textbox) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	_,err = buf.WriteString(html2.EscapeString(t.Text()))
	return
}

// Set the value of the text. Returns itself for chaining
func (t *Textbox) SetText(s string) page.ControlI {
	t.value = s
	t.SetAttribute("value", s)
	return t.this()
}

func (t *Textbox) Text() string {
	return t.value
}

func (t *Textbox) SetValue(v interface{}) page.ControlI {
	s := fmt.Sprintf("%T", v)
	t.this().SetText(s)
	return t.this()
}

func (t *Textbox) Value() interface{} {
	return t.this().Text()
}

func (t *Textbox) SetMaxLength(len int) *MaxLengthValidator {
	t.maxLength = len
	v := MaxLengthValidator{Length:len}
	t.ValidateWith(v)
	return &v
}

func (t *Textbox) MaxLength() int {
	return t.maxLength
}

func (t *Textbox) SetMinLength(len int) *MinLengthValidator {
	if len <= 0 {
		panic ("Cannot set minimum length to zero or less.")
	}
	t.minLength = len
	v := MinLengthValidator{Length:len}
	t.ValidateWith(v)
	return &v
}

func (t *Textbox) MinLength() int {
	return t.minLength
}

func (t *Textbox) SetPlaceholder(s string) TextboxI {
	t.SetAttribute("placeholder", s)
	return t.this()
}

func (t *Textbox) Placeholder() string {
	return t.Attribute("placeholder")
}

// SetType sets the type of textbox this is. Pass it a TEXTBOX_TYPE... constant normally, thought you can pass
// any string and it will become the input type
func (t *Textbox) SetType(typ string) TextboxI {
	t.typ = typ
	t.Refresh() // can't change this without completely redrawing the control
	return t.this()
}

// SetColumns sets the visible width of the text control. Each table is an approximate with of a character, and is browser
// dependent, so its not a very good way of setting the width. The css width property is more accurate. Also, this is
// only the visible width, not the maximum number of characters.
func (t *Textbox) SetColumns(columns int) {
	t.columns = columns
	if columns <= 0 {
		panic("Invalid table value.")
	}
	t.Refresh()
}

// SetRows sets the number of rows the Textbox will have. A value of 0 produces an input tag, and a value of 1 or greater produces a textarea tag.
func (t *Textbox) SetRows(rows int)  {
	if rows < 0 {
		panic("Invalid row value.")
	}
	if rows == 0 {
		t.Tag = "input"
		t.IsVoidTag = true
	} else {
		t.Tag = "textarea"
		t.IsVoidTag = false
	}
	t.rows = rows
	t.Refresh()
}

func (t *Textbox) SetReadOnly(r bool) {
	t.readonly = r
	t.AddRenderScript("attr", "readonly", "")
}

func (t *Textbox) SetSanitizer(s Sanitizer) {
	t.sanitizer = s
}
func (t *Textbox) sanitize(s string) string {
	if t.sanitizer == nil {
		panic ("You have to create a sanitizer. Not having a sanitizer is too dangerous.")
	}
	return t.sanitizer.Sanitize(s)
}

// Validate will first check for the Required attribute, and if set, will make sure a value is in the text field. It
// will then check the validators in the order assigned. The first invalid value found will return false.
func (t *Textbox) Validate() bool {
	text := t.Text()
	if t.Required() && text == "" {
		if t.ErrorForRequired == "" {
			t.SetValidationError(t.T("A value is required"))
		} else {
			t.SetValidationError(t.ErrorForRequired)
		}
		return false
	}

	if t.validators != nil {
		for _,v := range t.validators {
			if msg := v.Validate(t.Page().GoraddTranslator(), t.value); msg != "" {
				t.SetValidationError(msg)
				return false
			}
		}
	}
	return true
}

// UpdateFormValues is an internal function that lets us reflect the value of the textbox on the web page
func (t *Textbox) UpdateFormValues(ctx *page.Context) {
	id := t.Id()

	if v,ok := ctx.FormValue(id); ok {
		t.value = t.sanitize(v)
	}
}

// MarshalState is an internal function to save the state of the control
func (t *Textbox) MarshalState(m types.MapI) {
	m.Set("text", t.Text())
}

// UnmarshalState is an internal function to restore the state of the control
func (t *Textbox) UnmarshalState(m types.MapI) {
	if m.Has("text") {
		s,_ := m.GetString("text")
		t.value = s
	}
}

//TODO
func (t *Textbox) Serialize(buf []byte) {

}

func (t *Textbox) Unserialize(data interface{}) {

}


type MinLengthValidator struct {
	Length int
	Message string
}

func (v MinLengthValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if len(s) < v.Length {
		if v.Message == "" {
			return fmt.Sprintf (t.Translate("Enter at least %d characters"), v.Length)
		} else {
			return v.Message
		}
	}
	return
}

type MaxLengthValidator struct {
	Length int
	Message string
}

func (v MaxLengthValidator) Validate(t page.Translater, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if len(s) > v.Length {
		if v.Message == "" {
			return fmt.Sprintf (t.Translate("Enter at most %d characters"), v.Length)
		} else {
			return v.Message
		}
	}
	return
}