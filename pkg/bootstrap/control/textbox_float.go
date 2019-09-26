package control

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type FloatTextboxI interface {
	control.FloatTextboxI
}

type FloatTextbox struct {
	control.FloatTextbox
}

func NewFloatTextbox(parent page.ControlI, id string) *FloatTextbox {
	t := new(FloatTextbox)
	t.Init(t, parent, id)
	return t
}

func (t *FloatTextbox) ΩDrawingAttributes(ctx context.Context) html.Attributes {
	a := t.FloatTextbox.ΩDrawingAttributes(ctx)
	a.AddClass("form-control")
	return a
}

// Use FloatTextboxCreator to create a textbox that only accepts numbers.
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
	MinValue *control.FloatLimit
	// MaxValue is the maximum value the user can enter. If the user enter more
	// than this amount, or enters something that is not an integer, it will fail validation
	// and the FormFieldWrapper will show an error.
	MaxValue *control.FloatLimit
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
	// Reuse subclass
	sub := control.FloatTextboxCreator{
		Placeholder:    c.Placeholder,
		Type:           c.Type,
		MinLength:      c.MinLength,
		MaxLength:      c.MaxLength,
		RowCount:       c.RowCount,
		ReadOnly:       c.ReadOnly,
		ControlOptions: c.ControlOptions,
		SaveState:      c.SaveState,
		MinValue: c.MinValue,
		MaxValue: c.MaxValue,
		Value: c.Value,
	}
	sub.Init(ctx, ctrl)
}

// GetFloatTextbox is a convenience method to return the control with the given id from the page.
func GetFloatTextbox(c page.ControlI, id string) *FloatTextbox {
	return c.Page().GetControl(id).(*FloatTextbox)
}


func init() {
	gob.RegisterName("bootstrap.floattextbox", new(FloatTextbox))
}
