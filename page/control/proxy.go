package control

import (
	localPage "goradd/page"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/page/event"
    "bytes"
    "context"
    "github.com/spekary/goradd/javascript"
    "github.com/spekary/goradd/html"
    html2 "html"
	"fmt"
)


type Proxy struct {
	localPage.Control
}

func NewProxy(parent page.ControlI) *Proxy {
	p := Proxy{}
	p.Init(parent)
	return &p
}

func (p *Proxy) Init(parent page.ControlI) {
	p.Control.Init(p, parent)
    p.SetShouldAutoRender(true)
    p.SetActionValue(javascript.JsCode(`$j(this).data("grAv")`))
}

// OnClick is a shortcut for adding a click event handler that is particular to buttons. It debounces the click, to
// prevent potential accidental multiple form submissions.
func (p *Proxy) OnClick(actions... action.ActionI) {
	p.On(event.Click().Terminating().Delay(5).Blocking(), actions...)
}

func (p *Proxy) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
    response := p.Form().Response()
    p.This().PutCustomScript(ctx, response)
    p.GetActionScripts(response)
    return
}

// DrawAsLink draws the proxy as a link. Generally, only do this if you are actually linking to a page. If not, use
// a button.
func (p *Proxy) DrawAsLink(ctx context.Context,
    label string,
    actionValue string,
    attributes *html.Attributes,
    buf *bytes.Buffer,
) (err error) {
	if attributes == nil {
		attributes = html.NewAttributes()
	}
	if !attributes.Has("href") {
		attributes.Set("href", "javascript:;")
	}
	return p.DrawAsTag(ctx, label, actionValue, attributes, "a", false, buf)
}

func (p *Proxy) DrawAsTag(ctx context.Context,
    label string,
    actionValue string,
    attributes *html.Attributes,
    tag string,
    dontEscape bool,
    buf *bytes.Buffer,
) (err error) {
    a := html.NewAttributes()
    a.SetDataAttribute("grProxy", p.Id())

    if actionValue != "" {
        a.SetDataAttribute("grAv", actionValue)
    }

    if attributes != nil {
		a.Merge(attributes) // will only apply defaults that are not in attributes
	}

    if !dontEscape {
        label = html2.EscapeString(label)
    }

    _, err = buf.WriteString(html.RenderTag(tag, a, label))
    return
}

func (p *Proxy) DrawAsButton(ctx context.Context,
    label string,
    actionValue string,
    attributes *html.Attributes,
    dontEscape bool,
    buf *bytes.Buffer,
) (err error) {
    a := html.NewAttributes()
    a.Set("onclick", "return false")
    a.Set("type", "button")
    if attributes != nil {
		a.Merge(attributes)
	}
    p.DrawAsTag(ctx, label, actionValue, a, "button", dontEscape, buf)
    return
}

// Attributes returns attributes that can be included in any tag to attach a proxy to the tag.
func (p *Proxy) Attributes(actionValue string) *html.Attributes {
    a := html.NewAttributes()
    a.SetDataAttribute("grProxy", p.Id())

    if actionValue != "" {
        a.SetDataAttribute("grAv", actionValue)
    }

    return a
}

// WrapEvent is an internal function to allow the control to customize its treatment of event processing.
func (p *Proxy) WrapEvent(eventName string, selector string, eventJs string) string {
	return fmt.Sprintf(`$j('#%s').on('%s', '[data-gr-proxy="%s"]', function(event, ui){%s});`, p.Form().Id(), eventName, p.Id(), eventJs)
}

