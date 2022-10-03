package control

import (
	"context"
	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
	"io"
)

const DropdownSelect = "gr-bs-dropdownselect"

type DropdownI interface {
	control.UnorderedListI
	SetAsNavItem(bool) DropdownI
	ButtonAttributes() html5tag.Attributes
	MenuAttributes() html5tag.Attributes
	OnClick(action action.ActionI) DropdownI
}

// A Dropdown implements the Bootstrap Dropdown control.
// The Dropdown is a button combined with a list. The button shows the list, allowing the user
// to select an item from the list.
type Dropdown struct {
	control.UnorderedList
	asNavItem        bool
	buttonAttributes html5tag.Attributes
	menuAttributes   html5tag.Attributes
}

func NewDropdown(parent page.ControlI, id string) *Dropdown {
	l := &Dropdown{}
	l.Self = l
	l.Init(parent, id)
	return l
}

func (l *Dropdown) Init(parent page.ControlI, id string) {
	l.UnorderedList.Init(parent, id)
	l.Tag = "div"
	l.buttonAttributes = html5tag.NewAttributes()
	l.menuAttributes = html5tag.NewAttributes()

	// Trigger a DropdownSelect whenever an anchor with an href of "#" is clicked.
	// EventValue will be the value of the list item.
	l.On(event.Click().Selector(`a[href="#"][class~="dropdown-item"]`).Capture(),
		action.Trigger(l.ID(), DropdownSelect, javascript.JsCode(`g$(event.target).data("grEv")`)))
}

// this() supports object oriented features by giving easy access to the virtual function interface.
func (l *Dropdown) this() DropdownI {
	return l.Self.(DropdownI)
}

func (l *Dropdown) SetAsNavItem(asNavItem bool) DropdownI {
	l.asNavItem = asNavItem
	return l.this()
}

// ButtonAttributes returns the attributes for the internal button
// of the dropdown. If you change them, be sure to call Refresh().
func (l *Dropdown) ButtonAttributes() html5tag.Attributes {
	return l.buttonAttributes
}

// MenuAttributes returns the attributes for the internal button
// of the dropdown. If you change them, be sure to call Refresh().
func (l *Dropdown) MenuAttributes() html5tag.Attributes {
	return l.menuAttributes
}

// OnClick sets the action to take when a link in the Dropdown is selected.
// It will only respond to links whose href is "#", which indicates its an empty link.
// The ActionValue will be the id of the link clicked.
func (l *Dropdown) OnClick(a action.ActionI) DropdownI {
	l.On(DropdownSelectEvent(), a)
	return l.this()
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *Dropdown) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.UnorderedList.DrawingAttributes(ctx)
	a.AddClass("dropdown")
	a.Set("grctl", "bs-dropdown")
	if l.asNavItem {
		a.AddClass("nav-item")
	}
	return a
}

func (l *Dropdown) DrawInnerHtml(_ context.Context, w io.Writer) {
	var h string
	btnAttr := l.buttonAttributes.Copy()
	btnAttr.AddClass("dropdown-toggle").
		Set("href", "#").
		Set("role", "button").
		SetData("bsToggle", "dropdown").
		Set("aria-expanded", "false")
	if l.asNavItem {
		h = html5tag.RenderTag("a",
			btnAttr.AddClass("nav-link"),
			l.Text())
	} else {
		h = html5tag.RenderTag("a",
			btnAttr.AddClass("btn"),
			l.Text())
	}
	hItems := l.this().GetItemsHtml(l.ListItems())

	menuAttr := l.menuAttributes.Copy()
	menuAttr.AddClass("dropdown-menu")
	h += html5tag.RenderTag("ul", menuAttr, hItems)
	page.WriteString(w, h)
	return
}

func (l *Dropdown) GetItemsHtml(items []*control.ListItem) string {
	// make sure the list items have the correct classes before drawing them
	for _, item := range items {
		if item.Anchor() == "" {
			item.SetAnchor("#")
			item.AnchorAttributes().SetData("grEv", item.Value())
		}
		item.AnchorAttributes().AddClass("dropdown-item")
	}
	return l.UnorderedList.GetItemsHtml(items)
}

// Serialize serializes the state of the control for the pagestate
func (l *Dropdown) Serialize(e page.Encoder) {
	l.UnorderedList.Serialize(e)
	if err := e.Encode(l.asNavItem); err != nil {
		panic(err)
	}
	if err := e.Encode(l.buttonAttributes); err != nil {
		panic(err)
	}
	if err := e.Encode(l.menuAttributes); err != nil {
		panic(err)
	}
}

// Deserialize reconstructs the control from the page state.
func (l *Dropdown) Deserialize(d page.Decoder) {
	l.UnorderedList.Deserialize(d)

	if err := d.Decode(&l.asNavItem); err != nil {
		panic(err)
	}
	if err := d.Decode(&l.buttonAttributes); err != nil {
		panic(err)
	}
	if err := d.Decode(&l.menuAttributes); err != nil {
		panic(err)
	}
}

type DropdownCreator struct {
	// ID is the id attribute for the html object and the id for the page control
	ID string
	// Text is the label that will appear in the dropdown button
	Text string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []control.ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// AsNavItem determines whether to draw as an item in a navbar, or just a regular button.
	AsNavItem bool
	// ButtonAttributes are additional attributes that will be assigned to the button.
	ButtonAttributes html5tag.Attributes
	// MenuAttributes are additional attributes that will be assigned to the menu.
	MenuAttributes html5tag.Attributes
	// OnClick is the action to take when a link is clicked or selected. It will only respond
	// to anchor tags whose href is set to "#". The EventValue will be the value of the item clicked.
	OnClick action.ActionI
	// ControlOptions are additional settings for the control.
	// If this is part of a Navbar, you should add the "nav-item" class.
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c DropdownCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewDropdown(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c DropdownCreator) Init(ctx context.Context, ctrl DropdownI) {
	sub := control.UnorderedListCreator{
		Items:          c.Items,
		DataProvider:   c.DataProvider,
		DataProviderID: c.DataProviderID,
		ControlOptions: c.ControlOptions,
	}

	sub.Init(ctx, ctrl)
	ctrl.SetText(c.Text)
	ctrl.SetAsNavItem(c.AsNavItem)
	if c.ButtonAttributes != nil {
		ctrl.ButtonAttributes().Merge(c.ButtonAttributes)
	}
	if c.MenuAttributes != nil {
		ctrl.MenuAttributes().Merge(c.MenuAttributes)
	}

	if c.OnClick != nil {
		ctrl.OnClick(c.OnClick)
	}
}

// GetDropdown is a convenience method to return the control with the given id from the page.
func GetDropdown(c page.ControlI, id string) *Dropdown {
	return c.Page().GetControl(id).(*Dropdown)
}

func init() {
	page.RegisterControl(&Dropdown{})
}

func DropdownSelectEvent() *page.Event {
	return page.NewEvent(DropdownSelect)
}
