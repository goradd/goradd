package control

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control/data"
	page2 "goradd-project/override/page"
	"reflect"
)

// SelectList is typically a dropdown list with a single selection. Items are selected by id number, and the SelectList
// completely controls the ids in the list. Create the list by calling AddItem or AddItems to add ListItemI objects.
// Or, use the embedded DataManager to load items. Set the size attribute if you want to display it as a
// scrolling list rather than a dropdown list.
type SelectList struct {
	page2.Control
	ItemList
	data.DataManager
	selectedId string
}

func NewSelectList(parent page.ControlI, id string) *SelectList {
	t := &SelectList{}
	t.Init(t, parent, id)
	return t
}

func (l *SelectList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.Control.Init(self, parent, id)
	l.ItemList = NewItemList(l)
	l.Tag = "select"
}

func (l *SelectList) Validate(ctx context.Context) bool {
	if v := l.Control.Validate(ctx); !v {
		return false
	}

	if l.IsRequired() && l.selectedId == "" {
		if l.ErrorForRequired == "" {
			l.SetValidationError(l.T("A selection is required"))
		} else {
			l.SetValidationError(l.ErrorForRequired)
		}
		return false
	}
	return true
}

// UpdateFormValues is an internal function that lets us reflect the value of the selection on the web override
func (l *SelectList) UpdateFormValues(ctx *page.Context) {
	id := l.ID()

	if v, ok := ctx.FormValue(id); ok {
		l.selectedId = v
	}
}

// SelectedItem will return the currently selected item. If no item has been selected, it will return the first item
// in the list, since that is what will be showing in the selection list, and will update its internal pointer to
// make the first item the current selection.
func (l *SelectList) SelectedItem() ListItemI {
	if l.selectedId == "" {
		if l.Len() == 0 {
			return nil
		} else {
			l.selectedId = l.items[0].ID()
			return l.items[0]
		}
	}
	return l.GetItem(l.selectedId)
}

// SetSelectedId sets the current selection to the given id. You must ensure that the item with the id exists, it will
// not attempt to make sure the item exists.
func (l *SelectList) SetSelectedID(id string) {
	l.selectedId = id
	l.AddRenderScript("val", id)
}

// Value implements the Valuer interface for general purpose value getting and setting
func (l *SelectList) Value() interface{} {
	if i := l.SelectedItem(); i == nil {
		return nil
	} else {
		return i.Value()
	}
}

// SetValue implements the Valuer interface for general purpose value getting and setting
func (l *SelectList) SetValue(v interface{}) {
	s := fmt.Sprintf("%v", v)
	id, _ := l.GetItemByValue(s)
	l.SetSelectedID(id)
}

func (l *SelectList) IntValue() int {
	if i := l.SelectedItem(); i == nil {
		return 0
	} else {
		return i.IntValue()
	}
}

func (l *SelectList) StringValue() string {
	if i := l.SelectedItem(); i == nil {
		return ""
	} else {
		return i.StringValue()
	}
}

func (l *SelectList) SelectedLabel() string {
	item := l.SelectedItem()
	if item != nil {
		return item.Label()
	}
	return ""
}

// MarshalState is an internal function to save the state of the control
func (l *SelectList) MarshalState(m maps.Setter) {
	m.Set("sel", l.selectedId)
}

// UnmarshalState is an internal function to restore the state of the control
func (l *SelectList) UnmarshalState(m maps.Loader) {
	if v,ok := m.Load("sel"); ok {
		if s, ok := v.(string); ok {
			l.selectedId = s
		}
	}
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *SelectList) DrawingAttributes() *html.Attributes {
	a := l.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "selectlist")
	a.Set("name", l.ID()) // needed for posts
	if l.IsRequired() {
		a.Set("required", "")
	}
	return a
}

func (l *SelectList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	if l.HasDataProvider() {
		l.GetData(ctx, l)
	}
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *SelectList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			tag := "optgroup"
			innerhtml := l.getItemsHtml(item.ListItems())
			attributes := item.Attributes().Clone()
			attributes.Set("label", item.Label())
			h += html.RenderTag(tag, attributes, innerhtml)
		} else {
			attributes := item.Attributes().Clone()
			attributes.Set("value", item.ID())
			if l.selectedId == item.ID() {
				attributes.Set("selected", "")
			}

			h += html.RenderTag("option", attributes, item.RenderLabel())
		}
	}
	return h
}

// SetData overrides the default data setter to add objects to the item list. The result is kept in memory currently.
func (l *SelectList) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if kind != reflect.Slice {
		panic("you must call SetData with a slice")
	}

	l.ItemList.Clear()
	l.AddListItems(data)
}
