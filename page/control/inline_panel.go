package control

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/control_base"
)

// InlinePanel is a Goradd control that is a basic "span" tag.
type InlinePanel struct {
	control_base.Panel
}

func NewInlinePanel(parent page.ControlI) *InlinePanel {
	p := &InlinePanel{}
	p.Tag = "span"
	p.Init(p, parent)
	return p
}
