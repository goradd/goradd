package control

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"strconv"
)

type FloatTextboxI interface {
	TextboxI
	SetMinValue(minValue float64, invalidMessage string) FloatTextboxI
	SetMaxValue(maxValue float64, invalidMessage string) FloatTextboxI
}

// FloatTextbox is a text control that ensures a valid floating point number is entered in the field.
type FloatTextbox struct {
	Textbox
}

func NewFloatTextbox(parent page.ControlI, id string) *FloatTextbox {
	t := &FloatTextbox{}
	t.Self = t
	t.Init(parent, id)
	return t
}

func (t *FloatTextbox) Init(parent page.ControlI, id string) {
	t.Textbox.Init(parent, id)
	t.ValidateWith(FloatValidator{})
}

func (t *FloatTextbox) this() FloatTextboxI {
	return t.Self.(FloatTextboxI)
}

func (t *FloatTextbox) SetMinValue(minValue float64, invalidMessage string) FloatTextboxI {
	t.ValidateWith(MinFloatValidator{minValue, invalidMessage})
	return t.this()
}

func (t *FloatTextbox) SetMaxValue(maxValue float64, invalidMessage string) FloatTextboxI {
	t.ValidateWith(MaxFloatValidator{maxValue, invalidMessage})
	return t.this()
}

func (t *FloatTextbox) Value() interface{} {
	return t.Float64()
}

func (t *FloatTextbox) Float64() float64 {
	text := t.Textbox.Text()
	v, _ := strconv.ParseFloat(text, 64)
	return v
}

func (t *FloatTextbox) Float32() float32 {
	text := t.Textbox.Text()
	v, _ := strconv.ParseFloat(text, 32)
	return float32(v)
}

func (t *FloatTextbox) SetFloat64(v float64) *FloatTextbox {
	t.Textbox.SetValue(v)
	return t
}

func (t *FloatTextbox) SetFloat32(v float32) *FloatTextbox {
	t.Textbox.SetValue(v)
	return t
}

func (t *FloatTextbox) SetValue(v interface{}) page.ControlI {
	t.Textbox.SetValue(v)
	return t.this()
}

type FloatValidator struct {
	Message string
}

func (v FloatValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		if msg == "" {
			return c.GT("Please enter a number.")
		} else {
			return v.Message
		}
	}
	return
}

type MinFloatValidator struct {
	MinValue float64
	Message  string
}

func (v MinFloatValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.ParseFloat(s, 64); val < v.MinValue {
		if msg == "" {
			return fmt.Sprintf(c.GT("Enter at least %f"), v.MinValue)
		} else {
			return v.Message
		}
	}
	return
}

type MaxFloatValidator struct {
	MaxValue float64
	Message  string
}

func (v MaxFloatValidator) Validate(c page.ControlI, s string) (msg string) {
	if s == "" {
		return "" // empty textbox is checked elsewhere
	}
	if val, _ := strconv.ParseFloat(s, 64); val < v.MaxValue {
		if msg == "" {
			return fmt.Sprintf(c.GT("Enter at most %f"), v.MaxValue)
		} else {
			return v.Message
		}
	}
	return
}

type FloatLimit struct {
	Value          float64
	InvalidMessage string
}

// FloatTextboxCreator creates a textbox that only accepts numbers.
// Pass it to AddControls of a control, or as a Child of
// a FormFieldWrapper.
type FloatTextboxCreator struct {
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
	MinValue *FloatLimit
	// MaxValue is the maximum value the user can enter. If the user enter more
	// than this amount, or enters something that is not an integer, it will fail validation
	// and the FormFieldWrapper will show an error.
	MaxValue *FloatLimit
	// Value is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Value interface{}

	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c FloatTextboxCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewFloatTextbox(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Textboxes to initialize a control with the
// creator. You do not normally need to call this.
func (c FloatTextboxCreator) Init(ctx context.Context, ctrl FloatTextboxI) {
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

// GetFloatTextbox is a convenience method to return the control with the given id from the page.
func GetFloatTextbox(c page.ControlI, id string) *FloatTextbox {
	return c.Page().GetControl(id).(*FloatTextbox)
}

func init() {
	gob.Register(MaxFloatValidator{})
	gob.Register(MinFloatValidator{})
	gob.Register(FloatValidator{})
	page.RegisterControl(&FloatTextbox{})
}
