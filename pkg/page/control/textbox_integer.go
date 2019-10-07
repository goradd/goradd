package control

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"strconv"
)

type IntegerTextboxI interface {
	TextboxI
	SetMinValue(minValue int, invalidMessage string) IntegerTextboxI
	SetMaxValue(maxValue int, invalidMessage string) IntegerTextboxI
}

type IntegerTextbox struct {
	Textbox
	minValue *int
	maxValue *int
}

func NewIntegerTextbox(parent page.ControlI, id string) *IntegerTextbox {
	t := &IntegerTextbox{}
	t.Init(t, parent, id)
	return t
}

func (t *IntegerTextbox) Init(self TextboxI, parent page.ControlI, id string) {
	t.Textbox.Init(self, parent, id)
	t.ValidateWith(IntValidator{})
}

func (t *IntegerTextbox) this() IntegerTextboxI {
	return t.Self.(IntegerTextboxI)
}

// SetMinValue creates a validator that makes sure the value of the text box is at least the
// given value. Specify your own error message, or leave the error message blank and a standard error message will
// be presented if the value is not valid.
func (t *IntegerTextbox) SetMinValue(minValue int, invalidMessage string) IntegerTextboxI {
	t.ValidateWith(MinIntValidator{minValue, invalidMessage})
	t.minValue = new(int)
	*t.minValue = minValue
	return t.this()
}

// SetMaxValue creates a validator that makes sure the value of the text box is at most the
// given value. Specify your own error message, or leave the error message blank and a standard error message will
// be presented if the value is not valid.
func (t *IntegerTextbox) SetMaxValue(maxValue int, invalidMessage string) IntegerTextboxI {
	t.ValidateWith(MaxIntValidator{maxValue, invalidMessage})
	t.maxValue = new(int)
	*t.maxValue = maxValue
	return t.this()
}

func (t *IntegerTextbox) SetValue(v interface{}) page.ControlI {
	t.Textbox.SetValue(v)
	newValue := t.Int()
	if t.minValue != nil && *t.minValue > newValue {
		panic("Setting IntegerTextbox to a value less than minimum value.")
	}
	if t.maxValue != nil && *t.maxValue < newValue {
		panic("Setting IntegerTextbox to a value greater than the maximum value.")
	}
	return t.this()
}

func (t *IntegerTextbox) SetInt(v int) IntegerTextboxI {
	t.Textbox.SetValue(v)
	if t.minValue != nil && *t.minValue > v {
		panic("Setting IntegerTextbox to a value less than minimum value.")
	}
	if t.maxValue != nil && *t.maxValue < v {
		panic("Setting IntegerTextbox to a value greater than the maximum value.")
	}
	return t.this()
}

func (t *IntegerTextbox) Value() interface{} {
	return t.Int()
}

func (t *IntegerTextbox) Int() int {
	text := t.Textbox.Text()
	v, _ := strconv.Atoi(text)
	return v
}

func (t *IntegerTextbox) Int64() int64 {
	text := t.Textbox.Text()
	i64, _ := strconv.ParseInt(text, 10, 0)
	return i64
}

type IntValidator struct {
	Message string
}

func (v IntValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _, err := strconv.Atoi(s); err != nil {
		if v.Message == "" {
			return c.T("Please enter an integer.")
		} else {
			return v.Message
		}
	}
	return
}

type MinIntValidator struct {
	MinValue int
	Message  string
}

func (v MinIntValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.Atoi(s); val < v.MinValue {
		if v.Message == "" {
			return fmt.Sprintf(c.ΩT("Enter at least %d"), v.MinValue)
		} else {
			return v.Message
		}
	}
	return
}

type MaxIntValidator struct {
	MaxValue int
	Message  string
}

func (v MaxIntValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.Atoi(s); val < v.MaxValue {
		if v.Message == "" {
			return fmt.Sprintf(c.ΩT("Enter at most %d"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}

type IntegerLimit struct {
	Value          int
	InvalidMessage string
}

// Use IntegerTextboxCreator to create an integer textbox.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type IntegerTextboxCreator struct {
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
	// MinValue is the minimum value the user can enter. If the user does not
	// enter at least this amount, or enters something that is not an integer, it will fail validation
	// and the FormFieldWrapper will show an error.
	MinValue *IntegerLimit
	// MaxValue is the maximum value the user can enter. If the user enter more
	// than this amount, or enters something that is not an integer, it will fail validation
	// and the FormFieldWrapper will show an error.
	MaxValue *IntegerLimit
	// Value is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Value interface{}

	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c IntegerTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewIntegerTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Textboxes to initialize a control with the
// creator. You do not normally need to call this.
func (c IntegerTextboxCreator) Init(ctx context.Context, ctrl IntegerTextboxI) {
	if c.MinValue != nil {
		ctrl.SetMinValue(c.MinValue.Value, c.MinValue.InvalidMessage)
	}
	if c.MaxValue != nil {
		ctrl.SetMaxValue(c.MaxValue.Value, c.MaxValue.InvalidMessage)
	}
	if c.Value != nil {
		ctrl.SetValue(c.Value)
	}
	// Reuse subclass
	sub := TextboxCreator{
		Placeholder:    c.Placeholder,
		Type:           c.Type,
		MinLength:      c.MinLength,
		MaxLength:      c.MaxLength,
		RowCount:       c.RowCount,
		ReadOnly:       c.ReadOnly,
		ControlOptions: c.ControlOptions,
		SaveState:      c.SaveState,
	}
	sub.Init(ctx, ctrl)
}

// GetIntegerTextbox is a convenience method to return the control with the given id from the page.
func GetIntegerTextbox(c page.ControlI, id string) *IntegerTextbox {
	return c.Page().GetControl(id).(*IntegerTextbox)
}

func init() {
	gob.Register(MaxIntValidator{})
	gob.Register(MinIntValidator{})
	gob.Register(IntValidator{})
	gob.Register(IntegerTextbox{})
}
