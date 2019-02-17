package control

import (
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"bytes"
	"github.com/goradd/goradd/pkg/html"
	"fmt"
	"context"
	html2 "html"
)

type ProxyI interface {
	page.ControlI
	LinkHtml(label string,
		actionValue string,
		attributes *html.Attributes,
	) string
	TagHtml(label string,
		actionValue string,
		attributes *html.Attributes,
		tag string,
		rawHtml bool,
	) string
	ButtonHtml(label string,
		eventActionValue string,
		attributes *html.Attributes,
		rawHtml bool,
	) string
	OnSubmit(actions ...action.ActionI) page.EventI
}

type Proxy struct {
	page.Control
}

func NewProxy(parent page.ControlI) *Proxy {
	p := Proxy{}
	p.Init(parent)
	return &p
}

func (p *Proxy) Init(parent page.ControlI) {
	p.Control.Init(p, parent, "")
	p.SetShouldAutoRender(true)
	p.SetActionValue(javascript.JsCode(`$j(this).data("grAv")`))
}

func (p *Proxy) this() ProxyI {
	return p.Self.(ProxyI)
}

// OnSubmit is a shortcut for adding a click event handler that is particular to buttons. It debounces the click, to
// prevent potential accidental multiple form submissions. All events fired after this event fires will be lost. It is
// intended to be used when the action will result in navigating to a new page.
func (p *Proxy) OnSubmit(actions ...action.ActionI) page.EventI {
	return p.On(event.Click().Terminating().Delay(250), actions...)
}

// Draw is used by the form engine to draw the control. As a proxy, there is no html to draw, but this is where the scripts attached to the
// proxy get sent to the response. This should get drawn by the auto-drawing routine, since proxies are not rendered in templates.
func (p *Proxy) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	response := p.ParentForm().Response()
	// p.this().ΩPutCustomScript(ctx, response) // Proxies should not have custom scripts?
	p.GetActionScripts(response)
	p.ΩPostRender(ctx, buf)
	return
}

// LinkHtml renders the proxy as a link. Generally, only do this if you are actually linking to a page. If not, use
// a button.
func (p *Proxy) LinkHtml(label string,
	actionValue string,
	attributes *html.Attributes,
) string {
	if attributes == nil {
		attributes = html.NewAttributes()
	}
	if !attributes.Has("href") {
		attributes.Set("href", "javascript:;")
	}
	return p.TagHtml(label, actionValue, attributes, "a", false)
}

// TagHtml lets you customize the tag that will be used to embed the proxy.
func (p *Proxy) TagHtml(label string,
	actionValue string,
	attributes *html.Attributes,
	tag string,
	rawHtml bool,
) string {
	a := html.NewAttributes()
	a.SetDataAttribute("grProxy", p.ID())

	if actionValue != "" {
		a.SetDataAttribute("grAv", actionValue)
	}

	if attributes != nil {
		a.Merge(attributes) // will only apply defaults that are not in attributes
	}

	if !rawHtml {
		label = html2.EscapeString(label)
	}

	return html.RenderTagNoSpace(tag, a, label)
}

// ButtonHtml outputs the proxy as a button tag.
// eventActionValue becomes the event's value parameter
func (p *Proxy) ButtonHtml(label string,
	eventActionValue string,
	attributes *html.Attributes,
	rawHtml bool,
) string {
	a := html.NewAttributes()
	a.Set("onclick", "return false")
	a.Set("type", "button")
	if attributes != nil {
		a.Merge(attributes)
	}
	return p.TagHtml(label, eventActionValue, a, "button", rawHtml)
}

// Attributes returns attributes that can be included in any tag to attach a proxy to the tag.
func (p *Proxy) ActionAttributes(actionValue string) *html.Attributes {
	a := html.NewAttributes()
	a.SetDataAttribute("grProxy", p.ID())

	if actionValue != "" {
		a.SetDataAttribute("grAv", actionValue)
	}

	return a
}

// WrapEvent is an internal function to allow the control to customize its treatment of event processing.
func (p *Proxy) WrapEvent(eventName string, selector string, eventJs string) string {
	return fmt.Sprintf(`$j('#%s').on('%s', '[data-gr-proxy="%s"]', function(event, ui){%s});`, p.ParentForm().ID(), eventName, p.ID(), eventJs)
}
