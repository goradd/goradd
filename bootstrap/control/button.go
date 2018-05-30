package control

import (
	grctl "github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"goradd/app"
)
type Button struct {
	grctl.Button
	style ButtonStyle
	size ButtonSize
}

const ButtonClass = "btn"

type ButtonStyle string

const (
	ButtonStylePrimary ButtonStyle = "btn-primary"
	ButtonStyleSecondary = "btn-secondary"
	ButtonStyleSuccess = "btn-success"
	ButtonStyleInfo = "btn-info"
	ButtonStyleWarning = "btn-warning"
	ButtonStyleDanger = "btn-danger"
	ButtonStyleLight = "btn-light"
	ButtonStyleDark = "btn-dark"
	ButtonStyleLink = "btn-link"

	ButtonStyleOutlinePrimary = "btn-outline-primary"
	ButtonStyleOutlineSecondary = "btn-outline-secondary"
	ButtonStyleOutlineSuccess = "btn-outline-success"
	ButtonStyleOutlineInfo = "btn-outline-info"
	ButtonStyleOutlineWarning = "btn-outline-warning"
	ButtonStyleOutlineDanger = "btn-outline-danger"
	ButtonStyleOutlineLight = "btn-outline-light"
	ButtonStyleOutlineDark = "btn-outline-dark"
)

type ButtonSize string

const (
	ButtonSizeLarge ButtonSize = "btn-lg"
	ButtonSizeMedium= ""
	ButtonSizeSmall= "btn-sm"
)

// Add ButtonBlock as a class to a button to make it span a full block
const ButtonBlock = "btn-block"

// Creates a new standard html button
func NewButton(parent page.ControlI) *Button {
	b := &Button{}
	b.Init(b, parent)
	return b
}

func (b *Button) Init(self page.ControlI, parent page.ControlI) {
	b.Button.Init(self, parent)
	b.style = ButtonStyleSecondary // default
	app.LoadBootstrap(b.Form())
}

// SetButtonStyle will set the button's style to one of the predefined bootstrap styles.
func (b *Button) SetButtonStyle (style ButtonStyle) {
	b.style = style
}

// SetButtonsSize sets the size class of the button.
func (b *Button) SetButtonSize (size ButtonSize) {
	b.size = size
}

func (b *Button) DrawingAttributes() *html.Attributes {
	a := b.Button.DrawingAttributes()
	a.AddClass(ButtonClass)
	a.AddClass(string(b.style))
	a.AddClass(string(b.size))
	return a
}
