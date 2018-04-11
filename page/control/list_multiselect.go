package control

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/html"
	"bytes"
	"context"
	"strings"
	"strconv"
)

// MultiselectList is a generic list box which allows multiple selections. It is here for completeness, but is not used
// very often since it doesn't present an intuitive interface and is very browser dependent on what is presented.
// A Checkboxlist is better.
type MultiselectList struct {
	page.Control
	ItemList
	selectedIds map[string]bool
}

func NewMultiselectList(parent page.ControlI) *MultiselectList {
	l := &MultiselectList{}
	l.ItemList = NewItemList(l)
	l.Init(l, parent)
	l.selectedIds = map[string]bool{}
	return l
}


func (l *MultiselectList) Init(self page.ControlI, parent page.ControlI) {
	l.Control.Init(self, parent)
	l.Tag = "select"
}

func (l *MultiselectList) SetSize(size int) *MultiselectList {
	l.SetAttribute("size", strconv.Itoa(size))
	l.Refresh()
	return l
}

func (l *MultiselectList) Size() int  {
	a := l.Attribute("size")
	if a == "" {
		return 0
	} else {
		s,err := strconv.Atoi(a)
		if err != nil {
			return 0
		}
		return s
	}
}


func (l *MultiselectList) Validate() bool {
	if l.Required() && len(l.selectedIds) == 0 {
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
func (l *MultiselectList) UpdateFormValues(ctx *page.Context) {
	id := l.Id()

	l.selectedIds = map[string]bool{}
	if v,ok := ctx.FormValue(id); ok {
		a := strings.Split(v, ",")
		for _,v := range a {
			l.selectedIds[v] = true
		}
	}
}

func (l *MultiselectList) SelectedItems() []ListItemI {
	items := []ListItemI{}
	if len(l.selectedIds) == 0 {
		return nil
	}
	for id := range l.selectedIds {
		item := l.FindById(id)
		if item != nil {
			items = append(items, item)
		}
	}
	return items
}

// SetSelectedIds sets the current selection to the given ids. You must ensure that the items with the ids exist, it will
// not attempt to make sure the items exist.
func (l *MultiselectList) SetSelectedIds(ids []string) {
	if ids == nil {
		return
	}

	for _,id := range ids {
		l.selectedIds[id] = true
	}
	l.Refresh()
}

// Value implements the Valuer interface for general purpose value getting and setting
func (l *MultiselectList) Value() interface{} {
	return l.selectedIds
}

// SetValue implements the Valuer interface for general purpose value getting and setting
func (l *MultiselectList) SetValue(v interface{})  {
	l.selectedIds = map[string]bool{}
	switch ids := v.(type) {
	case string:
		a := strings.Split(ids, ",")
		for _,v := range a {
			l.selectedIds[v] = true
		}

	case []string:
		for _,v := range ids {
			l.selectedIds[v] = true
		}

	case ListItemI:
		l.selectedIds[ids.Id()] = true

	case []ListItemI:
		for _,v := range ids {
			l.selectedIds[v.Id()] = true
		}

	default:
		panic ("Unknown id list type")
	}
}

func (l *MultiselectList) SelectedLabels() []string {
	labels := []string{}

	for id := range l.selectedIds {
		item := l.FindById(id)
		if item != nil {
			labels = append(labels, item.Label())
		}
	}
	return labels
}

// MarshalState is an internal function to save the state of the control
func (l *MultiselectList) MarshalState(m types.MapI) {
	var ids = []string{}
	for id := range l.selectedIds {
		ids = append(ids, id)
	}
	m.Set("sel", ids)
}

// UnmarshalState is an internal function to restore the state of the control
func (l *MultiselectList) UnmarshalState(m types.MapI) {
	l.selectedIds = map[string]bool{}

	if m.Has("sel") {
		s := m.Get("sel")

		if ids,ok := s.([]string); ok {
			for _,id := range ids {
				l.selectedIds[id] = true
			}
		}
	}
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *MultiselectList) DrawingAttributes() *html.Attributes {
	a := l.Control.DrawingAttributes()
	a.Set("name", l.Id())	// needed for posts
	a.Set("multiple", "")
	if l.Required() {
		a.Set("required", "")
	}
	return a
}

func (l *MultiselectList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *MultiselectList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _,item := range items {
		if item.HasChildItems() {
			tag := "optgroup"
			innerhtml := l.getItemsHtml(item.ListItems())
			attributes := item.Attributes().Clone()
			attributes.Set("label", item.Label())
			h += html.RenderTag(tag, attributes, innerhtml) + "\n"
		} else {
			attributes := item.Attributes().Clone()
			attributes.Set("value", item.Id())
			if _,ok := l.selectedIds[item.Id()]; ok {
				attributes.Set("selected", "")
			}
			h += html.RenderTag("option",attributes, item.Label()) + "\n"
		}
	}
	return h
}