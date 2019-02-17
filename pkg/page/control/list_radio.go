package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
)

type RadioListI interface {
	CheckboxListI
}

// RadioList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-table table, or a single-table
// scrolling list much like a standard html select list.
type RadioList struct {
	CheckboxList
}

// NewRadioList creates a new RadioList control.
func NewRadioList(parent page.ControlI, id string) *RadioList {
	l := &RadioList{}
	l.Init(l, parent, id)
	return l
}

// Init is called by subclasses.
func (l *RadioList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.CheckboxList.Init(self, parent, id)
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) ΩDrawingAttributes() *html.Attributes {
	a := l.CheckboxList.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "radiolist")
	return a
}

// ΩRenderItem is called by the framework to render a single item in the list.
func (l *RadioList) ΩRenderItem(tag string, item ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", l.ID())
	attributes.Set("value", item.ID())
	attributes.Set("type", "radio")
	if l.selectedIds[item.ID()] {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	h = html.RenderLabel(html.NewAttributes().Set("for", item.ID()), item.Label(), ctrl, l.labelDrawingMode)
	attributes = item.Attributes().Copy()
	attributes.SetID(item.ID() + "_cell")
	attributes.AddClass("gr-cbl-item")
	h = html.RenderTag(tag, attributes, h)
	return
}

// Value returns the single selected value of the list and satisfies the Valuer interface.
// It returns nil if no item is selected.
func (l *RadioList) Value() interface{} {
	a := l.SelectedValues()
	if len(a) == 0 {
		return nil
	} else {
		return a[0]
	}
}

// SelectedValue returns the single selected value of the list as a string.
func (l *RadioList) SelectedValue() string {
	return l.Value().(string)
}

// SelectedLabel returns the label of the currently selected item, or an empty string of no item is selected.
func (l *RadioList) SelectedLabel() string {
	a := l.SelectedLabels()
	if len(a) == 0 {
		return ""
	} else {
		return a[0]
	}
}

// SetValue sets the selection to the item corresponding to the given value, and satisfies the Valuer interface.
func (l *RadioList) SetValue(v interface{}) {
	l.SetSelectedValue(v)
}

// SetSelectedValue sets the selection to the item corresponding to the given value.
func (l *RadioList) SetSelectedValue(v interface{}) {
	if v == nil {
		l.SetSelectedID("")
		return
	}

	id, item := l.GetItemByValue(v)
	if item != nil {
		l.SetSelectedID(id)
	}
}

// SetSelectedID sets selection to the item whose id corresponds to the given value.
func (l *RadioList) SetSelectedID(id string) {
	if id == "" {
		l.selectedIds = map[string]bool{}
	} else {
		l.selectedIds = map[string]bool{id:true}
	}
	l.Refresh()
}


func (l *RadioList) ΩUpdateFormValues(ctx *page.Context) {
	controlID := l.ID()

	if ctx.RequestMode() == page.Ajax {
		if v, ok := ctx.CheckableValue(controlID); ok {
			if s, ok := v.(string); ok {
				l.selectedIds = map[string]bool{l.ID() + "_" + s: true}
			}
		}
	} else {
		if v, ok := ctx.FormValue(controlID); ok {
			l.selectedIds = map[string]bool{v:true}
		}
	}
}
