package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/crypt"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"html"
	"io"
	"strings"
)

type ProxyI interface {
	page.ControlI
	LinkHtml(label string,
		actionValue string,
		attributes html5tag.Attributes,
	) string
	TagHtml(label string,
		actionValue string,
		attributes html5tag.Attributes,
		tag string,
		rawHtml bool,
	) string
	ButtonHtml(label string,
		eventActionValue string,
		attributes html5tag.Attributes,
		rawHtml bool,
	) string
	OnSubmit(action action.ActionI) *page.Event
}

// Proxy is a control that attaches events to controls. It is useful for attaching
// similar events to a series of controls, like all the links in a table, or all the buttons in button bar.
// You can also use it to draw a series of links or buttons. The proxy differentiates between the different objects
// that are sending it events by the ActionValue that you gave the proxy.
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
	page.ControlBase
}

// NewProxy creates a new proxy. The parent must be the wrapping control of the objects that the proxy will manage.
func NewProxy(parent page.ControlI, id string) *Proxy {
	p := &Proxy{}
	p.Self = p
	p.Init(parent, id)
	return p
}

func (p *Proxy) Init(parent page.ControlI, id string) {
	p.ControlBase.Init(parent, id)
	p.SetShouldAutoRender(true)
	p.SetActionValue(javascript.JsCode(`goradd.proxyVal(event)`))
}

func (p *Proxy) this() ProxyI {
	return p.Self.(ProxyI)
}

// OnSubmit is a shortcut for adding a click event handler that is particular to buttons. It debounces the click, to
// prevent potential accidental multiple form submissions. All events fired after this event fires will be lost. It is
// intended to be used when the action will result in navigating to a new page.
func (p *Proxy) OnSubmit(action action.ActionI) page.ControlI {
	return p.On(event.Click().Terminating().Delay(250), action)
}

// Draw is used by the form engine to draw the control. As a proxy, there is no html to draw, but this is where the scripts attached to the
// proxy get sent to the response. This should get drawn by the auto-drawing routine, since proxies are not rendered in templates.
func (p *Proxy) Draw(ctx context.Context, w io.Writer) {
	response := p.ParentForm().Response()
	// p.this().PutCustomScript(ctx, response) // Proxies should not have custom scripts?

	p.GetActionScripts(response)
	p.PostRender(ctx, w)
	return
}

// LinkHtml renders the proxy as a link. To conform to the html standard and accessibility guidelines,
// links should only be used to navigate away from the page, so the action of your proxy should lead to
// that kind of behavior. Otherwise, use ButtonHtml.
func (p *Proxy) LinkHtml(ctx context.Context,
	label string,
	actionValue string,
	attributes html5tag.Attributes,
) string {
	if attributes == nil {
		attributes = html5tag.NewAttributes()
	}
	attributes.Set("onclick", "return false;") // make sure we do not follow the link if javascript is on.
	var href string
	if attributes.Has("href") {
		href = attributes.Get("href")
	} else {
		href = page.GetContext(ctx).HttpContext.URL.RequestURI() // for non-javascript compatibility
		if offset := strings.Index(href, page.HtmlVarAction); offset >= 0 {
			href = href[:offset-1] // remove the variables we placed here ourselves
		}
	}

	// These next two lines allow the proxy to work even when javascript is off.
	av := page.HtmlVarAction + "=" + p.ID() + "_" + actionValue
	av += "&" + page.HtmlVarPagestate + "=" + crypt.SessionEncryptUrlValue(ctx, p.Page().StateID())

	if !strings.ContainsRune(href, '?') {
		href += "?" + av
	} else {
		href += "&" + av
	}
	attributes.Set("href", href)
	return p.TagHtml(label, actionValue, attributes, "a", false)
}

// TagHtml lets you customize the tag that will be used to embed the proxy.
func (p *Proxy) TagHtml(label string,
	actionValue string,
	attributes html5tag.Attributes,
	tag string,
	labelIsHtml bool,
) string {
	a := html5tag.NewAttributes()
	a.SetData("grProxy", p.ID())

	if actionValue != "" {
		a.SetData("grAv", actionValue)
	}

	if attributes != nil {
		a.Merge(attributes) // will only apply defaults that are not in attributes
	}

	if !labelIsHtml {
		label = html.EscapeString(label)
	}

	return html5tag.RenderTagNoSpace(tag, a, label)
}

// ButtonHtml outputs the proxy as a button tag.
// actionValue becomes the event's ControlValue parameter
func (p *Proxy) ButtonHtml(label string,
	actionValue string,
	attributes html5tag.Attributes,
	labelIsHtml bool,
) string {
	a := html5tag.NewAttributes()
	a.Set("onclick", "return false")  // To prevent a return from activating the button
	a.Set("type", "submit")           // To support non-javascript situations
	a.Set("name", page.HtmlVarAction) // needed for non-javascript posts
	buttonValue := p.ID() + "_" + actionValue
	a.Set("value", buttonValue) // needed for non-javascript posts

	if attributes != nil {
		a.Merge(attributes)
	}

	// TODO: We can possibly do actionValue differently now since its already in the value above
	return p.TagHtml(label, actionValue, a, "button", labelIsHtml)
}

// ActionAttributes returns attributes that can be included in any tag to attach a proxy to the tag.
func (p *Proxy) ActionAttributes(actionValue string) html5tag.Attributes {
	a := html5tag.NewAttributes()
	a.SetData("grProxy", p.ID())

	if actionValue != "" {
		a.SetData("grAv", actionValue)
	}

	return a
}

// WrapEvent is an internal function to allow the control to customize its treatment of event processing.
func (p *Proxy) WrapEvent(eventName string, _ string, eventJs string, options map[string]interface{}) string {
	// This attaches the event to the parent control.
	return fmt.Sprintf(`g$('%s').on('%s', '[data-gr-proxy="%s"]', function(event, eventData){%s}, %s);`, p.Parent().ID(), eventName, p.ID(), eventJs, javascript.ToJavaScript(options))
}

type On struct {
	Event *page.Event
	Action action.ActionI
}

type ProxyCreator struct {
	// ID is the id of the proxy. Proxies do not draw, so this id will not show up in the html, but you can
	// use it to get the proxy from the page.
	ID string
	// On is a shortcut to assign a single action to an event. If you want a proxy that responds to more than
	// one event or action, use On in the ControlOptions struct
	On On
	page.ControlOptions
}

func (c ProxyCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewProxy(parent, c.ID)
	if c.On.Event != nil {
		ctrl.On(c.On.Event, c.On.Action)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	return ctrl
}

// GetProxy is a convenience method to return the button with the given id from the page.
func GetProxy(c page.ControlI, id string) *Proxy {
	return c.Page().GetControl(id).(*Proxy)
}

func init() {
	page.RegisterControl(&Proxy{})
}