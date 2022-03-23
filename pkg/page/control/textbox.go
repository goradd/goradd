package control

import (
	"context"
	"encoding/gob"
	"fmt"
	"html"
	"io"
	"strconv"

	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

const (
	TextboxTypeDefault  = "text"
	TextboxTypePassword = "password"
	TextboxTypeSearch   = "search"
	TextboxTypeNumber   = "number" // Puts little arrows in box, will need to widen it.
	TextboxTypeEmail    = "email"  // see TextEmail. Prevents submission of RFC5322 email addresses (Gogh Fir <gf@example.com>)
	TextboxTypeTel      = "tel"
	TextboxTypeUrl      = "url"
)

// A Validater can be added to a textbox to validate its input on the server side.
// A textbox can have more than one validater.
// A number of built-in validators are provided.
type Validater interface {
	// Validate evaluates the input, and returns an empty string if the input is valid, and an error string to display
	// to the user if the input does not pass the validator.
	Validate(page.ControlI, string) string
}

type TextboxI interface {
	page.ControlI
	SetType(typ string) TextboxI
	Sanitize(string) string
	SetPlaceholder(s string) TextboxI
	SetMaxLength(len int) *MaxLengthValidator
	SetMinLength(len int) *MinLengthValidator
	SetRowCount(rows int) TextboxI
	SetColumnCount(columns int) TextboxI
	SetReadOnly(r bool) TextboxI
	SetValue(interface{}) page.ControlI
}

// Textbox is a goradd control that outputs an "input" html tag with a "type" attribute
// of "text", or one of the text-like types, like "password", "search", etc.
type Textbox struct {
	page.ControlBase

	typ string

	validators []Validater

	minLength int
	maxLength int

	value string

	columnCount int
	rowCount    int

	readonly bool
}

// NewTextbox creates a new goradd textbox html widget.
func NewTextbox(parent page.ControlI, id string) *Textbox {
	t := &Textbox{}
	t.Self = t
	t.Init(parent, id)
	return t
}

// Init initializes a textbox. Normally you will not call this directly.
func (t *Textbox) Init(parent page.ControlI, id string) {
	t.ControlBase.Init(parent, id)

	t.Tag = "input"
	t.IsVoidTag = true
	t.typ = "text" // default
	t.SetHasNoSpace(true)
}

func (t *Textbox) this() TextboxI {
	return t.Self.(TextboxI)
}

// ValidateWith adds a Validater to the validator list.
func (t *Textbox) ValidateWith(v Validater) {
	t.validators = append(t.validators, v)
}

// ResetValidators removes all validators
func (t *Textbox) ResetValidators() {
	t.validators = nil
}

// DrawingAttributes is called by the framework to retrieve the tag's private attributes at draw time.
func (t *Textbox) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := t.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "textbox")
	a.Set("name", t.ID()) // needed for posts
	if t.IsRequired() {
		a.Set("required", "")
	}
	if t.maxLength != 0 {
		a.Set("maxlength", strconv.Itoa(t.maxLength))
	}
	if t.rowCount == 0 { // single-line textbox
		a.Set("type", t.typ)
		a.Set("value", t.value)
		if t.columnCount != 0 {
			a.Set("size", strconv.Itoa(t.columnCount))
		}
	} else {
		a.Set("rows", strconv.Itoa(t.rowCount))
		if t.columnCount != 0 {
			a.Set("cols", strconv.Itoa(t.columnCount))
		}
	}
	a.AddValues("aria-labelledby", t.ID()) // spec says inputs should label themselves so screen reader will read out content of the input
	if t.readonly {
		a.Set("readonly", "")
	}
	return a
}

// DrawInnerHtml is an internal function that renders the inner html of a tag. In this case, it is rendering the inner
// text of a textarea
func (t *Textbox) DrawInnerHtml(_ context.Context, w io.Writer) {
	page.WriteString(w, html.EscapeString(t.Text()))
	return
}

// SetText sets the value of the text. Returns itself for chaining.
func (t *Textbox) SetText(s string) page.ControlI {
	if t.value != s {
		t.value = s
		t.AddRenderScript("val", s)
	}
	return t.this()
}

// Text returns the text entered by the user.
func (t *Textbox) Text() string {
	return t.value
}

// SetValue sets the text in the textbox. This satisfies the Valuer interface.
func (t *Textbox) SetValue(v interface{}) page.ControlI {
	s := fmt.Sprint(v)
	t.this().SetText(s)
	return t.this()
}

// Value returns the user entered text in the textbox.
func (t *Textbox) Value() interface{} {
	return t.this().Text()
}

// SetMaxLength sets the maximum length allowed in the textbox. The text will be limited by the
// browser, but the server side will also make sure that the text is not too big.
func (t *Textbox) SetMaxLength(len int) *MaxLengthValidator {
	t.maxLength = len
	v := MaxLengthValidator{Length: len}
	t.ValidateWith(v)
	return &v
}

// MaxLength returns the current maximum length setting.
func (t *Textbox) MaxLength() int {
	return t.maxLength
}

// SetMinLength will set the minimum length permitted. If the user does not enter enough text,
// an error message will be displayed upon submission of the form.
func (t *Textbox) SetMinLength(len int) *MinLengthValidator {
	if len < 0 {
		panic("Cannot set minimum length to less than zero.")
	}
	t.minLength = len
	v := MinLengthValidator{Length: len}
	t.ValidateWith(v)
	return &v
}

// MinLength returns the minimum length setting.
func (t *Textbox) MinLength() int {
	return t.minLength
}

// SetPlaceholder will set the html placeholder attribute, which puts text in the textbox
// when the textbox is empty as a hint to the user of what to enter.
func (t *Textbox) SetPlaceholder(s string) TextboxI {
	t.SetAttribute("placeholder", s)
	return t.this()
}

// Placeholder returns the value of the placeholder.
func (t *Textbox) Placeholder() string {
	return t.Attribute("placeholder")
}

// SetType sets the type of textbox this is. Pass it a TextboxType* constant normally,
// though you can pass any string and it will become the input type
func (t *Textbox) SetType(typ string) TextboxI {
	t.typ = typ
	t.Refresh() // can't change this without completely redrawing the control
	return t.this()
}

// SetColumnCount sets the visible width of the text control. Each table is an approximate with of
// a character, and is browser
// dependent, so its not a very good way of setting the width.
// The css width property is more accurate. Also, this is
// only the visible width, not the maximum number of characters.
func (t *Textbox) SetColumnCount(columns int) TextboxI {
	t.columnCount = columns
	if columns <= 0 {
		panic("Invalid table value.")
	}
	t.Refresh()
	return t.this()
}

// SetRowCount sets the number of rowCount the Textbox will have.
// A value of 0 produces an input tag, and a value of 1 or greater produces a textarea tag.
func (t *Textbox) SetRowCount(rows int) TextboxI {
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
	t.rowCount = rows
	t.Refresh()
	return t.this()
}

// SetReadOnly will disable editing by setting a browser attribute.
func (t *Textbox) SetReadOnly(r bool) TextboxI {
	t.readonly = r
	t.AddRenderScript("attr", "readonly", "")
	return t.this()
}

// Sanitize is called by the framework when taking in user input and strips it of potential
// malicious XSS scripts.
//
// The default uses a global sanitizer created at bootup.
// Override Sanitize in a subclass if you want a per-textbox sanitizer.
// This is a very difficult thing to get right, and depends a bit on your application on just
// how much you want to remove.
func (t *Textbox) Sanitize(s string) string {
	if config.GlobalSanitizer == nil {
		return s
	}
	return config.GlobalSanitizer.Sanitize(s)
}

// Validate will first check for the IsRequired attribute, and if set, will make sure a value is in the text field. It
// will then check the validators in the order assigned. The first invalid value found will return false.
func (t *Textbox) Validate(ctx context.Context) bool {
	if v := t.ControlBase.Validate(ctx); !v {
		return false
	}
	text := t.Text()
	if t.IsRequired() && text == "" {
		if t.ErrorForRequired == "" {
			t.SetValidationError(t.GT("A value is required"))
		} else {
			t.SetValidationError(t.ErrorForRequired)
		}
		return false
	}

	if t.validators != nil {
		for _, v := range t.validators {
			if msg := v.Validate(t.this(), t.value); msg != "" {
				t.SetValidationError(msg)
				return false
			}
		}
	}
	return true
}

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (t *Textbox) UpdateFormValues(ctx context.Context) {
	if t.readonly {
		// This would happen if someone was attempting to hack the browser.
		return
	}

	id := t.ID()

	if v, ok := page.GetContext(ctx).FormValue(id); ok {
		t.value = t.this().Sanitize(v)
	}
}

// MarshalState is an internal function to save the state of the control
func (t *Textbox) MarshalState(m page.SavedState) {
	m.Set("text", t.Text())
}

// UnmarshalState is an internal function to restore the state of the control
func (t *Textbox) UnmarshalState(m page.SavedState) {
	if v, ok := m.Load("text"); ok {
		if s, ok2 := v.(string); ok2 {
			t.value = s
		}
	}
}

type encodedTextbox struct {
	Typ         string
	Validators  []Validater
	MinLength   int
	MaxLength   int
	Value       string
	ColumnCount int
	RowCount    int
	Readonly    bool
}

// Serialize is used by the framework to serialize the textbox into the pagestate.
func (t *Textbox) Serialize(e page.Encoder) {
	t.ControlBase.Serialize(e)

	s := encodedTextbox{
		Typ:         t.typ,
		Validators:  t.validators,
		MinLength:   t.minLength,
		MaxLength:   t.maxLength,
		Value:       t.value,
		ColumnCount: t.columnCount,
		RowCount:    t.rowCount,
		Readonly:    t.readonly,
	}

	if err := e.Encode(s); err != nil {
		panic(err)
	}
}

// Deserialize is used by the pagestate serializer.
func (t *Textbox) Deserialize(d page.Decoder) {
	t.ControlBase.Deserialize(d)

	s := encodedTextbox{}

	if err := d.Decode(&s); err != nil {
		panic(err)
	}

	t.typ = s.Typ
	t.validators = s.Validators
	t.minLength = s.MinLength
	t.maxLength = s.MaxLength
	t.value = s.Value
	t.columnCount = s.ColumnCount
	t.rowCount = s.RowCount
	t.readonly = s.Readonly
}

// MinLengthValidator is a validator that checks that the user has entered a minimum length.
// It is set up automatically by calling SetMinValue.
type MinLengthValidator struct {
	Length  int
	Message string
}

// Validate runs the Validate logic to validate the control value.
func (v MinLengthValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if len(s) < v.Length {
		if v.Message == "" {
			return fmt.Sprintf(c.GT("Enter at least %d characters"), v.Length) // not a great translation, probably should be an Sprintf implementation
		} else {
			return v.Message
		}
	}
	return
}

// MaxLengthValidator is a Validater to test that the user did not enter too many characters.
type MaxLengthValidator struct {
	Length  int
	Message string
}

// Validate runs the Validate logic to validate the control value.
func (v MaxLengthValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if len(s) > v.Length {
		if v.Message == "" {
			return fmt.Sprintf(c.GT("Enter at most %d characters"), v.Length)
		} else {
			return v.Message
		}
	}
	return
}

// TextboxCreator creates a textbox. Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type TextboxCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	// Placeholder is the placeholder attribute of the textbox and shows as help text inside the field
	Placeholder string
	// Type is the type attribute of the textbox
	Type string
	// MinLength is the minimum number of characters that the user is required to enter. If the
	// length is less than this number, a validation error will be shown.
	MinLength int
	// MaxLength is the maximum number of characters that the user is required to enter. If the
	// length is more than this number, a validation error will be shown.
	MaxLength int
	// ColumnCount is the number of characters wide the textbox will be, and becomes the width attribute in the tag.
	// The actual width is browser dependent. For better control, use a width style property.
	ColumnCount int
	// RowCount creates a multi-line textarea with the given number of rows. By default the
	// textbox will expand vertically by this number of lines. Use a height style property for
	// better control of the height of a textbox.
	RowCount int
	// ReadOnly sets the readonly attribute of the textbox, which prevents it from being changed by the user.
	ReadOnly bool
	// SaveState will save the text in the textbox, to be restored if the user comes back to the page.
	// It is particularly helpful when the textbox is being used to filter the results of a query, so that
	// when the user comes back to the page, he does not have to type the filter text again.
	SaveState bool
	// Text is the initial value of the textbox. Generally you would not use this, but rather load the value in a separate Load step after creating the control.
	Text string

	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c TextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewTextbox(parent, c.ID)

	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Textboxes to initialize a control with the
// creator. You do not normally need to call this.
func (c TextboxCreator) Init(ctx context.Context, ctrl TextboxI) {
	if c.Placeholder != "" {
		ctrl.SetPlaceholder(c.Placeholder)
	}
	if c.Type != "" {
		ctrl.SetType(c.Type)
	}

	if c.MinLength != 0 {
		ctrl.SetMinLength(c.MinLength)
	}

	if c.MaxLength != 0 {
		ctrl.SetMaxLength(c.MaxLength)
	}
	if c.RowCount > 0 {
		ctrl.SetRowCount(c.RowCount)
	}
	if c.ColumnCount > 0 {
		ctrl.SetColumnCount(c.ColumnCount)
	}
	if c.ReadOnly {
		ctrl.SetReadOnly(true)
	}
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}

	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
}

// GetTextbox is a convenience method to return the control with the given id from the page.
func GetTextbox(c page.ControlI, id string) *Textbox {
	return c.Page().GetControl(id).(*Textbox)
}

func GetTextboxI(c page.ControlI, id string) TextboxI {
	return c.Page().GetControl(id).(TextboxI)
}

func init() {
	// gob.Register(&Textbox{}) register control.Textbox instead
	gob.Register(MaxLengthValidator{})
	gob.Register(MinLengthValidator{})
	page.RegisterControl(&Textbox{})
}
