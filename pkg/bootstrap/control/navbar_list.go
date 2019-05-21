package control

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/html"
	"bytes"
	"github.com/goradd/goradd/pkg/page/control"
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control/data"
)

type NavbarListI interface {
	page.ControlI
}

type NavbarList struct {
	page.Control
	control.ItemList
	subItemTag string
	data.DataManager
	Proxy *control.Proxy
}


func NavbarSelectEvent() page.EventI {
	e := &page.Event{JsEvent: "gr-bs-navbarselect"}
	e.ActionValue(javascript.JsCode("ui"))	// This will be the action value sent by the proxy...the id of the item
	return e
}

// TODO: Create a mechanism to post-process this event and have it automatically be loaded with the selected item

func NewNavbarList(parent page.ControlI, id string) *NavbarList {
	t := &NavbarList{}
	t.ItemList = control.NewItemList(t)
	t.Init(t, parent, id)
	return t
}

func (l *NavbarList) Init(self NavbarListI, parent page.ControlI, id string) {
	l.Control.Init(self, parent, id)
	l.Tag = "ul"
	l.subItemTag = "li"
	l.Proxy = control.NewProxy(l)

	l.Proxy.On(event.Click(),
		action.Trigger(l.ID(), "gr-bs-navbarselect", javascript.JsCode("g$(this).data('grAv')")))
}

func (l *NavbarList) this() NavbarListI {
	return l.Self.(NavbarListI)
}

func (l *NavbarList) ΩDrawTag(ctx context.Context) string {
	if l.DataManager.HasDataProvider() {
		l.GetData(ctx, l)
		defer l.Clear()
	}
	return l.Control.ΩDrawTag(ctx)
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *NavbarList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "navbarlist")
	a.AddClass("navbar-nav")
	return a
}

func (l *NavbarList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(ctx, l.ListItems(), false)
	buf.WriteString(h)
	return nil
}

func (l *NavbarList) getItemsHtml(ctx context.Context, items []control.ListItemI, hasParent bool) string {
	var h = ""

	for i, item := range items {
		if item.HasChildItems() {
			if hasParent {
				// A dropdown inside a dropdown
				h += fmt.Sprintf(
					`<a class="dropdown-item dropdown-toggle" id="%s_menu" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
        %s
    </a>
    <div class="dropdown-menu" aria-labelledby="%s_menu">`, item.ID(), item.RenderLabel(), item.ID())
				h += l.getItemsHtml(ctx, item.ListItems(), true)
				h += "</div>"
			} else {
				// top level menu
				var lastClass = ""

				if i == len(items) - 1 {
					// last item, so modify dropdown menu so it does not go off of screen
					// If there is only one item in the navbar, and this is the left navbar, this might cause a problem.
					// We can potentially fix that by asking the parent item if that is the situation.
					lastClass = "dropdown-menu-right"
				}
				h += fmt.Sprintf(
`<%s class="nav-item dropdown">
    <a class="nav-link dropdown-toggle" id="%s_menu" role="menu" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
        %s
    </a>
    <div class="dropdown-menu %s" aria-labelledby="%s_menu">`, l.subItemTag, item.ID(), item.RenderLabel(), lastClass, item.ID())
				h += l.getItemsHtml(ctx, item.ListItems(), true)
				h += fmt.Sprintf("</div></%s>", l.subItemTag)
			}
		} else {
			if item.IsDivider() {
				h += html.RenderTag("div", html.NewAttributes().AddClass("dropdown-divider"), "")
			} else if item.Disabled() {
				if !hasParent {
					h += fmt.Sprintf(`<li class="nav-item">
    <a class="nav-link disabled" href="#">%s</a>
</li>`, item.RenderLabel())

				} else {
					h += fmt.Sprintf(`<a class="dropdown-item disabled" href="#">%s</a>
</li>`, item.RenderLabel())
				}
			} else {
				itemH := item.RenderLabel()
				itemAttributes := item.Attributes().Copy()
				itemAttributes.AddClass("nav-item")
				linkAttributes := html.NewAttributes()


				if hasParent {
					itemAttributes.Set("role", "menuitem")
					linkAttributes.AddClass("dropdown-item")
				} else {
					linkAttributes.AddClass("nav-link")
				}

				if item.Anchor() == "" {
					itemH = l.Proxy.LinkHtml(ctx, itemH, item.ID(), linkAttributes)
				}
				if !hasParent {
					itemH = html.RenderTag(l.subItemTag, itemAttributes, itemH)
				}
				h += itemH
			}
		}
	}
	return h
}
