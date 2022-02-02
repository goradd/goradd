package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"io"
)

type NavbarListI interface {
	page.ControlI
	control.ItemListI
	control.DataManagerI
	OnSelect (action action.ActionI) page.ControlI
}

type NavbarList struct {
	page.ControlBase
	control.ItemList
	subItemTag string
	control.DataManager
}


func NewNavbarList(parent page.ControlI, id string) *NavbarList {
	t := &NavbarList{}
	t.Self = t
	t.ItemList = control.NewItemList(t)
	t.Init(parent, id)
	return t
}

func (l *NavbarList) Init(parent page.ControlI, id string) {
	l.ControlBase.Init(parent, id)
	l.Tag = "ul"
	l.subItemTag = "li"

	pxy := control.NewProxy(l, l.proxyID())

	pxy.On(event.Click(),
		action.Trigger(l.ID(), NavbarSelect, javascript.JsCode("g$(event.target).data('grAv')")))
	config.LoadBootstrap(l.ParentForm())
}

func (l *NavbarList) proxyID() string {
	return l.ID() + "-pxy"
}

func (l *NavbarList) ItemProxy() *control.Proxy {
	return control.GetProxy(l, l.proxyID())
}

func (l *NavbarList) this() NavbarListI {
	return l.Self.(NavbarListI)
}

func (l *NavbarList) DrawTag(ctx context.Context) string {
	if l.DataManager.HasDataProvider() {
		l.this().LoadData(ctx, l.this())
		defer l.ResetData() // prevent the data from being serialized and taking up space unnecessarily
	}
	return l.ControlBase.DrawTag(ctx)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *NavbarList) DrawingAttributes(ctx context.Context) html.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "navbarlist")
	a.AddClass("navbar-nav")
	return a
}

func (l *NavbarList) DrawInnerHtml(ctx context.Context, w io.Writer) (err error) {
	h := l.getItemsHtml(ctx, l.ListItems(), false)
	_,err = io.WriteString(w, h)
	return
}

func (l *NavbarList) getItemsHtml(ctx context.Context, items []*control.ListItem, hasParent bool) string {
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
				var itemClass string

				if i == len(items)-1 {
					// last item, so modify dropdown menu so it does not go off of screen
					// If there is only one item in the navbar, and this is the left navbar, this might cause a problem.
					// We can potentially fix that by asking the parent item if that is the situation.
					itemClass = "dropdown-menu-right "
				}
				// Let the item style it further
				itemClass += item.Attributes().Class()
				h += fmt.Sprintf(
					`<%s class="nav-item dropdown">
    <a class="nav-link dropdown-toggle" id="%s_menu" role="menu" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
        %s
    </a>
    <div class="dropdown-menu %s" aria-labelledby="%s_menu">`, l.subItemTag, item.ID(), item.RenderLabel(), itemClass, item.ID())
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
			} else if hasParent {
				itemH := item.RenderLabel()
				itemAttributes := item.Attributes().Copy()
				itemAttributes.AddClass("nav-item")
				linkAttributes := html.NewAttributes()
				itemAttributes.Set("role", "menuitem")
				linkAttributes.AddClass("dropdown-item")
				if !item.HasAnchor() {
					itemH = l.ItemProxy().LinkHtml(ctx, itemH, item.Value(), linkAttributes)
				}
				h += itemH
			} else {
				item.AnchorAttributes().AddClass("nav-link")
				itemH := item.RenderLabel()
				itemAttributes := item.Attributes().Copy()
				itemAttributes.AddClass("nav-item")

				if item.Anchor() == "" {
					linkAttributes := html.NewAttributes()
					linkAttributes.AddClass("nav-link")
					itemH = l.ItemProxy().LinkHtml(ctx, itemH, item.Value(), linkAttributes)
				}
				itemH = html.RenderTag(l.subItemTag, itemAttributes, itemH)
				h += itemH
			}
		}
	}
	return h
}

func (l *NavbarList) OnSelect (action action.ActionI) page.ControlI {
	return l.On(NavbarSelectEvent(), action)
}

func (l *NavbarList) Serialize(e page.Encoder) (err error) {
	if err = l.ControlBase.Serialize(e); err != nil {
		return
	}
	if err = l.ItemList.Serialize(e); err != nil {
		return
	}
	if err = l.DataManager.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(l.subItemTag); err != nil {
		return
	}
	return
}

func (l *NavbarList) Deserialize(dec page.Decoder) (err error) {
	if err = l.ControlBase.Deserialize(dec); err != nil {
		return
	}
	if err = l.ItemList.Deserialize(dec); err != nil {
		return
	}
	if err = l.DataManager.Deserialize(dec); err != nil {
		return
	}
	if err = dec.Decode(&l.subItemTag); err != nil {
		return
	}
	return
}

type NavbarListCreator struct {
	// ID is the control id of the html widget and must be unique to the page
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []interface{}
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	page.ControlOptions
	// OnSelect is the action to take when a list item is selected.
	OnSelect action.ActionI
}

func (c NavbarListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewNavbarList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c NavbarListCreator) Init(ctx context.Context, ctrl NavbarListI) {
	if c.Items != nil {
		ctrl.AddListItems(c.Items...)
	}

	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(control.DataBinder)
		ctrl.SetDataProvider(provider)
	}

	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.OnSelect != nil {
		ctrl.OnSelect(c.OnSelect)
	}
}


// GetNavbarList is a convenience method to return the control with the given id from the page.
func GetNavbarList(c page.ControlI, id string) *NavbarList {
	return c.Page().GetControl(id).(*NavbarList)
}

func init() {
	page.RegisterControl(&NavbarList{})
}

const NavbarSelect = "gr-bs-navbarselect"

func NavbarSelectEvent() *page.Event {
	return page.NewEvent(NavbarSelect)
}

