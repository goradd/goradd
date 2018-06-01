package control

import (
	"github.com/spekary/goradd/page"
	"fmt"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/html"
	"bytes"
	"context"
	html2 "html"
)

// SelectList is a typical dropdown list with a single selection. Items are selected by id number, and the SelectList
// completely controls the ids in the list. Create the list by calling AddItem or AddItems to add ListItemI objects.
type SelectList struct {
	page.Control
	ItemList
	selectedId string
}

func NewSelectList(parent page.ControlI) *SelectList {
	t := &SelectList{}
	t.ItemList = NewItemList(t)
	t.Init(t, parent)
	return t
}


func (l *SelectList) Init(self page.ControlI, parent page.ControlI) {
	l.Control.Init(self, parent)

	l.Tag = "select"
}

func (l *SelectList) Validate() bool {
	if v := l.Validate(); !v {
		return false
	}

	if l.Required() && l.selectedId == "" {
		if l.ErrorForRequired == "" {
			l.SetValidationError(l.T("A selection is required"))
		} else {
			l.SetValidationError(l.ErrorForRequired)
		}
		return false
	}
	return true
}

// UpdateFormValues is an internal function that lets us reflect the value of the selection on the web page
func (l *SelectList) UpdateFormValues(ctx *page.Context) {
	id := l.ID()

	if v,ok := ctx.FormValue(id); ok {
		l.selectedId = v
	}
}

func (l *SelectList) SelectedItem() ListItemI {
	if l.selectedId == "" {
		return nil
	}
	return l.FindByID(l.selectedId)
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
func (l *SelectList) SetValue(v interface{})  {
	s := fmt.Sprintf("%v")
	id, _ := l.FindByValue(s)
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
func (l *SelectList) MarshalState(m types.MapI) {
	m.Set("sel", l.selectedId)
}

// UnmarshalState is an internal function to restore the state of the control
func (l *SelectList) UnmarshalState(m types.MapI) {
	if m.Has("sel") {
		s,_ := m.GetString("sel")
		l.selectedId = s
	}
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *SelectList) DrawingAttributes() *html.Attributes {
	a := l.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "selectlist")
	a.Set("name", l.ID())	// needed for posts
	if l.Required() {
		a.Set("required", "")
	}
	return a
}

func (l *SelectList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *SelectList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _,item := range items {
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

			h += html.RenderTag("option",attributes, html2.EscapeString(item.Label()))
		}
	}
	return h
}