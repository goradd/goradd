package control

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
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

// Proxy is a control that attaches events to controls. It is useful for attaching
// similar events to a series of controls, like all the links in a table, or all the buttons in button bar.
// You can also use it to draw a series of links or buttons. The proxy differentiates between the different objects
// that are sending it events by the ActionValue that you given the proxy when it draws.
//
// To use a Proxy, create it in the control that wraps the controls the proxy will manage.
// Attach an event to the proxy control, and in the action handler, look for the ControlValue in the Action Value
// to know which of the controls sent the event. Draw the proxy with one of the following:
//   LinkHtml() - Output the proxy as a link
//   ButtonHtml() - Output the proxy as a button
//   TagHtml() - Output the proxy in any tag
//   ActionAttributes() - Returns attributes you can use in any custom control to attach a proxy
//
// The ProxyColumn of the Table object will use a proxy to draw items in a table column.
type Proxy struct {
	page.Control
}

// NewProxy creates a new proxy. The parent should be the wrapping control of the objects that the proxy will manage.
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
// eventActionValue becomes the event's ControlValue parameter
func (p *Proxy) ButtonHtml(label string,
	eventActionValue string,
	attributes *html.Attributes,
	rawHtml bool,
) string {
	a := html.NewAttributes()
	a.Set("onclick", "return false") // To prevent a return from activating the button
	a.Set("type", "submit") // To support non-javascript situations
	a.Set("name", page.HtmlVarAction) // needed for non-javascript posts
	buttonValue := p.ID() + "_" + eventActionValue
	a.Set("value", buttonValue) // needed for non-javascript posts

	if attributes != nil {
		a.Merge(attributes)
	}

	// TODO: We can possibly do actionValue differently now since its already in the value above
	return p.TagHtml(label, eventActionValue, a, "button", rawHtml)
}

// ActionAttributes returns attributes that can be included in any tag to attach a proxy to the tag.
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
	// This attaches the event to the parent control.
	return fmt.Sprintf(`$j('#%s').on('%s.grproxy', '[data-gr-proxy="%s"]', function(event, ui){%s});`, p.Parent().ID(), eventName, p.ID(), eventJs)
}
