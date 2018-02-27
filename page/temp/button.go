package temp
import (
	//"github.com/microcosm-cc/bluemonday"
	"net/mail"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	localPage "goradd/page"
)


type ButtonI interface {
	page.ControlI
}

type Button struct {
	localPage.Control
}

// Creates a new standard html button
func NewButton(parent page.ControlI, id string) *TextBox {
	b := &Button{}
	b.Init(b, parent, id)
	return t
}
