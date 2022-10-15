package control

import (
	"context"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
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
	p.Self = p
	p.Init(parent, id)
	return p
}

func (c *Span) Init(parent page.ControlI, id string) {
	c.Panel.Init(parent, id)
	c.Tag = "span"
}

func (c *Span) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := c.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "span")
	return a
}

// SpanCreator creates a "span" tag.
type SpanCreator struct {
	// ID is the HTML id of the tag.
	ID string
	// Text is the content inside the span tag.
	Text string
	// TextIsHtml specifies that the content is HTML and should not be escaped.
	TextIsHtml bool
	// Children are child controls of the span.
	Children []page.Creator
	// ControlOptions are standard options you can set on any GoRADD control.
	page.ControlOptions
}

func (c SpanCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewSpan(parent, c.ID)
	if c.Text != "" {
		ctrl.SetText(c.Text)
	}
	ctrl.SetTextIsHtml(c.TextIsHtml)
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	ctrl.AddControls(ctx, c.Children...)
	return ctrl
}

// GetSpan is a convenience method to return the button with the given id from the page.
func GetSpan(c page.ControlI, id string) *Span {
	return c.Page().GetControl(id).(*Span)
}

func init() {
	page.RegisterControl(&Span{})
}
