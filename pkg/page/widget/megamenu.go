package widget

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"io"
	"reflect"
	"strconv"
)

// MegaMenu is a control that is the basis for a navigation or action menu that is typically used
// in a large application. It implements the recommendations at https://www.levelaccess.com/challenges-mega-menus-standard-menus-make-accessible/
// to create a natively styled menu, which is keyboard navigable. It tries to avoid the aria-menu item, as
// that is particularly difficult to implement.
//
// If an item has sub items, it is treated simply as a menu opener, and clicking on that item will not
// produce an action. If a terminal item has an href, it will be output as a link. Otherwise, it will
// be output as a button and will fire the "MenuSelectEvent".
//
// There are so many ways to show menus. We attempt to provide some examples, but you might need
// to add additional styling for your situation.
type MegaMenu struct {
	page.ControlBase
	control.ItemList
	control.DataManager
}

type MegaMenuI interface {
	page.ControlI
	control.ItemListI
	control.DataManagerI
	GetItemsHtml(items []*control.ListItem, level int) string
	SetAriaLabel(l string) MegaMenuI
}


func NewMegaMenu(parent page.ControlI, id string) *MegaMenu {
	l := &MegaMenu{}
	l.Self = l
	l.Init(parent, id)
	pxy := control.NewProxy(l, l.ID() + "-pxy")
	pxy.On(event.Click(), action.Trigger(l.ID(), MegaMenuSelectEvent, javascript.JsCode{"g$(event.target).data('grAv')"}))

	return l
}

func (l *MegaMenu) Init(parent page.ControlI, id string) {
	l.ControlBase.Init(parent, id)
	l.ItemList = control.NewItemList(l)
	l.Tag = "nav"
}

// this() supports object oriented features by giving easy access to the virtual function interface.
func (l *MegaMenu) this() MegaMenuI {
	return l.Self.(MegaMenuI)
}


func (l *MegaMenu) getProxy() *control.Proxy {
	return control.GetProxy(l, l.ID()+"-pxy")
}


// SetAriaLabel sets the aria label that will be used in the tag. e.g. main menu
func (l *MegaMenu) SetAriaLabel(s string) MegaMenuI {
	l.SetAttribute("aria-label", s)
	return l.this()
}

func (l *MegaMenu) DrawTag(ctx context.Context) string {
	if l.HasDataProvider() {
		l.this().LoadData(ctx, l.this())
		defer l.ResetData()
	}
	return l.ControlBase.DrawTag(ctx)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *MegaMenu) DrawingAttributes(ctx context.Context) html.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "megamenu")
	return a
}

func (l *MegaMenu) DrawInnerHtml(ctx context.Context, w io.Writer) {
	h := l.this().GetItemsHtml(l.ListItems(), 1)
	h = html.RenderTag("ul", html.Attributes{"style":"list-style:none"}, h)
	page.WriteString(w, h)
	return
}

// GetItemsHtml is used by the framework to get the items for the html. It is exported so that
// it can be overridden by other implementations of an MegaMenu.
func (l *MegaMenu) GetItemsHtml(items []*control.ListItem, level int) string {
	var h = ""

	for _, item := range items {
		buttonId := l.ID() + "_" + item.Value()
		if item.HasChildItems() {
			innerhtml := l.this().GetItemsHtml(item.ListItems(), level + 1)
			innerhtml = html.RenderTag("ul", html.Attributes{"style":"list-style:none"}, innerhtml)
			innerhtml = html.RenderTag("div", html.Attributes{"role":"region", "aria-labeledby":buttonId}, innerhtml)
			buttonhtml := html.RenderTag("button", html.Attributes{"aria-expanded":`false`, "id":buttonId}, item.Label())
			innerhtml = html.RenderTag("div", html.Attributes{"role":"heading", "aria-level":strconv.Itoa(level+1)}, buttonhtml) + innerhtml
			h += html.RenderTag("li", item.Attributes(), innerhtml)
		} else {
			if item.HasAnchor() {
				a := html.RenderTag("a", item.Attributes(), item.RenderLabel())
				h += html.RenderTag("li", nil, a)
			} else {
				a := item.Attributes().Copy()
				a.SetID(buttonId)
				b := l.getProxy().ButtonHtml(item.Label(),
					item.Value(),
					a,
					false)
				h += html.RenderTag("li", nil, b)
			}
		}
	}
	return h
}

// SetData replaces the current list with the given data.
// ItemLister, ItemIDer, Labeler or Stringer types are accepted.
// This function can accept one or more lists of items, or
// single items. They will all get added to the top level of the list. To add sub items, get a list item
// and add items to it.
func (l *MegaMenu) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic("you must call SetData with a slice or array")
	}

	l.ItemList.Clear()
	l.AddListItems(data)
}

func (l *MegaMenu) Serialize(e page.Encoder) (err error) {
	if err = l.ControlBase.Serialize(e); err != nil {
		return
	}
	if err = l.ItemList.Serialize(e); err != nil {
		return
	}
	if err = l.DataManager.Serialize(e); err != nil {
		return
	}

	return
}

func (l *MegaMenu) Deserialize(dec page.Decoder) (err error) {
	if err = l.ControlBase.Deserialize(dec); err != nil {
		return
	}
	if err = l.ItemList.Deserialize(dec); err != nil {
		panic(err)
	}
	if err = l.DataManager.Deserialize(dec); err != nil {
		panic(err)
	}
	return
}


type MegaMenuCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []control.ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// AriaLabel is the aria label attribute.
	AriaLabel string
	OnMenuSelect action.ActionI
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c MegaMenuCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewMegaMenu(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c MegaMenuCreator) Init(ctx context.Context, ctrl MegaMenuI) {
	if c.Items != nil {
		ctrl.AddListItems(c.Items)
	}
	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(control.DataBinder)
		ctrl.SetDataProvider(provider)
	}
	if c.AriaLabel != "" {
		ctrl.SetAriaLabel(c.AriaLabel)
	}
	if c.OnMenuSelect != nil {
		ctrl.On(MegaMenuSelect(), c.OnMenuSelect)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
}

// GetMegaMenu is a convenience method to return the control with the given id from the page.
func GetMegaMenu(c page.ControlI, id string) *MegaMenu {
	return c.Page().GetControl(id).(*MegaMenu)
}

func init() {
	page.RegisterControl(&MegaMenu{})
}


const MegaMenuSelectEvent = "megamenuselect"

func MegaMenuSelect() *page.Event {
	return page.NewEvent(MegaMenuSelectEvent)
}

