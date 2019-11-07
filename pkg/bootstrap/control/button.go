package control

import (
	"context"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	grctl "github.com/goradd/goradd/pkg/page/control"
)

type ButtonI interface {
	grctl.ButtonI
	SetButtonStyle(style ButtonStyle) ButtonI
	SetButtonSize(size ButtonSize) ButtonI
	SetIsPrimary(isPrimary bool) ButtonI
}

type Button struct {
	grctl.Button
	style ButtonStyle
	size  ButtonSize
}

const ButtonClass = "btn"

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

type ButtonSize string

const (
	ButtonSizeLarge  ButtonSize = "btn-lg"
	ButtonSizeMedium            = ""
	ButtonSizeSmall             = "btn-sm"
)

// Add ButtonBlock as a class to a button to make it span a full block
const ButtonBlock = "btn-block"

// Creates a new standard html button
func NewButton(parent page.ControlI, id string) *Button {
	b := &Button{}
	b.Self = b
	b.Init(parent, id)
	return b
}

func (b *Button) Init(parent page.ControlI, id string) {
	b.Button.Init(parent, id)
	b.style = ButtonStyleSecondary // default
	config.LoadBootstrap(b.ParentForm())
}

func (b *Button) this() ButtonI {
	return b.Self.(ButtonI)
}

// SetButtonStyle will set the button's style to one of the predefined bootstrap styles.
func (b *Button) SetButtonStyle(style ButtonStyle) ButtonI {
	b.style = style
	return b.this()
}

// SetButtonsSize sets the size class of the button.
func (b *Button) SetButtonSize(size ButtonSize) ButtonI {
	b.size = size
	return b.this()
}

func (b *Button) DrawingAttributes(ctx context.Context) html.Attributes {
	a := b.Button.DrawingAttributes(ctx)
	a.AddClass(ButtonClass)
	a.AddClass(string(b.style))
	a.AddClass(string(b.size))
	return a
}

func (b *Button) SetIsPrimary(isPrimary bool) ButtonI {
	if isPrimary {
		b.style = ButtonStylePrimary
	} else {
		b.style = ButtonStyleSecondary
	}
	return b.this()
}

func (b *Button) Serialize(e page.Encoder) (err error) {
	if err = b.Button.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(b.style); err != nil {
		return err
	}

	if err = e.Encode(b.size); err != nil {
		return err
	}

	return
}


func (b *Button) Deserialize(d page.Decoder) (err error) {
	if err = b.Button.Deserialize(d); err != nil {
		return
	}

	if err = d.Decode(&b.style); err != nil {
		return
	}

	if err = d.Decode(&b.size); err != nil {
		return
	}

	return
}

// ButtonCreator is the initialization structure for declarative creation of buttons
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
	OnClick action.ActionI
	Style ButtonStyle
	Size ButtonSize
	IsPrimary bool
	ValidationType page.ValidationType
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
// creator. You do not normally need to call this.
func (c ButtonCreator) Init(ctx context.Context, ctrl ButtonI)  {
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
}