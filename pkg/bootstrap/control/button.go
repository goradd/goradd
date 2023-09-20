package control

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	grctl "github.com/goradd/goradd/pkg/page/control/button"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
)

// ButtonI is the Bootstrap control Button interface.
type ButtonI interface {
	grctl.ButtonI
	SetButtonStyle(style ButtonStyle) ButtonI
	SetButtonSize(size ButtonSize) ButtonI
}

// Button will draw a Bootstrap styled control.
type Button struct {
	grctl.Button
	style ButtonStyle
	size  ButtonSize
}

const ButtonClass = "btn"

// ButtonStyle represents a btn-* Bootstrap style class.
type ButtonStyle string

const (
	ButtonStylePrimary   ButtonStyle = "btn-primary"
	ButtonStyleSecondary             = "btn-secondary"
	ButtonStyleSuccess               = "btn-success"
	ButtonStyleInfo                  = "btn-info"
	ButtonStyleWarning               = "btn-warning"
	ButtonStyleDanger                = "btn-danger"
	ButtonStyleLight                 = "btn-light"
	ButtonStyleDark                  = "btn-dark"
	ButtonStyleLink                  = "btn-link"

	ButtonStyleOutlinePrimary   = "btn-outline-primary"
	ButtonStyleOutlineSecondary = "btn-outline-secondary"
	ButtonStyleOutlineSuccess   = "btn-outline-success"
	ButtonStyleOutlineInfo      = "btn-outline-info"
	ButtonStyleOutlineWarning   = "btn-outline-warning"
	ButtonStyleOutlineDanger    = "btn-outline-danger"
	ButtonStyleOutlineLight     = "btn-outline-light"
	ButtonStyleOutlineDark      = "btn-outline-dark"
)

// ButtonSize is one of the btn-* size classes.
type ButtonSize string

const (
	ButtonSizeLarge  ButtonSize = "btn-lg"
	ButtonSizeMedium            = ""
	ButtonSizeSmall             = "btn-sm"
)

// ButtonBlock is a class you can add to a button to make it span a full block
const ButtonBlock = "btn-block"

// NewButton creates a new bootstrap html button
func NewButton(parent page.ControlI, id string) *Button {
	b := new(Button)
	b.Init(b, parent, id)
	return b
}

// Init initializes the button
func (b *Button) Init(self any, parent page.ControlI, id string) {
	b.Button.Init(self, parent, id)
	b.style = ButtonStyleSecondary // default
	config.LoadBootstrap(b.ParentForm())
}

func (b *Button) this() ButtonI {
	return b.Self().(ButtonI)
}

// SetButtonStyle will set the button's style to one of the predefined bootstrap styles.
func (b *Button) SetButtonStyle(style ButtonStyle) ButtonI {
	b.style = style
	return b.this()
}

// SetButtonSize sets the size class of the button.
func (b *Button) SetButtonSize(size ButtonSize) ButtonI {
	b.size = size
	return b.this()
}

// DrawingAttributes returns the attributes to add to the tag just before the button is drawn.
func (b *Button) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := b.Button.DrawingAttributes(ctx)
	a.AddClass(ButtonClass)
	a.AddClass(string(b.style))
	a.AddClass(string(b.size))
	return a
}

// SetIsPrimary determines whether the button is styled as a primary or secondary button.
func (b *Button) SetIsPrimary(isPrimary bool) grctl.ButtonI {
	b.Button.SetIsPrimary(isPrimary)
	if isPrimary {
		b.style = ButtonStylePrimary
	} else {
		b.style = ButtonStyleSecondary
	}
	return b.this()
}

// Serialize serializes the state of the control for the pagestate
func (b *Button) Serialize(e page.Encoder) {
	b.Button.Serialize(e)

	if err := e.Encode(b.style); err != nil {
		panic(err)
	}

	if err := e.Encode(b.size); err != nil {
		panic(err)
	}

	return
}

// Deserialize reconstructs the control from the page state.
func (b *Button) Deserialize(d page.Decoder) {
	b.Button.Deserialize(d)

	if err := d.Decode(&b.style); err != nil {
		panic(err)
	}

	if err := d.Decode(&b.size); err != nil {
		panic(err)
	}
}

// ButtonCreator is the initialization structure for the declarative creation of the control.
type ButtonCreator struct {
	// ID is the control id
	ID string
	// Text is the text displayed in the button
	Text string
	// OnSubmit is the action to take when the button is submitted. Use this specifically
	// for buttons that move to other pages or processes transactions, as it debounces the button
	// and waits until all other actions complete
	OnSubmit action.ActionI
	// OnClick is an action to take when the button is pressed. Do not specify both
	// a OnSubmit and OnClick.
	OnClick        action.ActionI
	Style          ButtonStyle
	Size           ButtonSize
	IsPrimary      bool
	ValidationType event.ValidationType
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c ButtonCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewButton(parent, c.ID)

	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator.
func (c ButtonCreator) Init(ctx context.Context, ctrl ButtonI) {
	sub := grctl.ButtonCreator{
		ID:             c.ID,
		Text:           c.Text,
		OnSubmit:       c.OnSubmit,
		OnClick:        c.OnClick,
		ValidationType: c.ValidationType,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
	if c.Style != "" {
		ctrl.SetButtonStyle(c.Style)
	}
	if c.Size != "" {
		ctrl.SetButtonSize(c.Size)
	}
	if c.IsPrimary {
		ctrl.SetIsPrimary(true)
	}
}

// GetButton is a convenience method to return the button with the given id from the page.
func GetButton(c page.ControlI, id string) *Button {
	return c.Page().GetControl(id).(*Button)
}

func init() {
	page.RegisterControl(&Button{})
	var s ButtonStyle
	gob.Register(s)
}
