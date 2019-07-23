package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

type SpanI interface {
	PanelI
}

// Span is a Goradd control that is a basic "span" wrapper. Use it to style and listen to events on a span. It
// can also be used as the basis for more advanced javascript controls.
type Span struct {
	Panel
}

func NewSpan(parent page.ControlI, id string) *Span {
	p := &Span{}
	p.Init(p, parent, id)
	return p
}

func (c *Span) Init(self SpanI, parent page.ControlI, id string) {
	c.Panel.Init(self, parent, id)
	c.Tag = "span"
}

func (c *Span) ΩDrawingAttributes() *html.Attributes {
	a := c.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "span")
	return a
}

type SpanCreator struct {
	ID string
	Text string
	TextIsHtml bool
	Children []page.Creator
	page.ControlOptions
}

func (c SpanCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewSpan(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	ctrl.SetTextIsHtml(c.TextIsHtml)
	ctrl.ApplyOptions(c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
	return ctrl
}
