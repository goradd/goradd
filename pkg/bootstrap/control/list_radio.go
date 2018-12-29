package control

import (
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/html"
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

func NewRadioList(parent page.ControlI, id string) *RadioList {
	l := &RadioList{}
	l.Init(l, parent, id)
	return l
}

func (l *RadioList) Init(self RadioListI, parent page.ControlI, id string) {
	l.CheckboxList.Init(self, parent, id)
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

func (l *RadioList) renderItem(item control.ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", l.ID())
	attributes.Set("value", item.ID())
	attributes.Set("type", "radio")
	if l.IsIdSelected(item.ID()) {
		attributes.Set("checked", "")
	}
	attributes.AddClass("form-check-input")
	ctrl := html.RenderVoidTag("input", attributes)

	h = html.RenderLabel(html.NewAttributes().Set("for", item.ID()).AddClass("form-check-label"), item.Label(), ctrl, html.LabelAfter)
	attributes = item.Attributes().Clone()
	attributes.AddClass("form-check")
	if l.isInline {
		attributes.AddClass("form-check-inline")
	}
	h = html.RenderTag("div", attributes, h)
	return
}

func (l *RadioList) Value() interface{} {
	a := l.SelectedValues()
	if len(a) == 0 {
		return nil
	} else {
		return a[0]
	}
}

func (l *RadioList) SelectedValue() string {
	return l.Value().(string)
}

func (l *RadioList) SelectedLabel() string {
	a := l.SelectedLabels()
	if len(a) == 0 {
		return ""
	} else {
		return a[0]
	}
}

func (l *RadioList) SetValue(v interface{}) {
	l.SetSelectedValue(v)
}

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

func (l *RadioList) SetSelectedID(id string) {
	l.SetSelectedIds([]string{id})
	if id == "" {
		l.SetSelectedIds(nil)
	} else {
		l.SetSelectedIds([]string{id})
	}
}


func (l *RadioList) UpdateFormValues(ctx *page.Context) {
	controlID := l.ID()

	if ctx.RequestMode() == page.Ajax {
		if v, ok := ctx.CheckableValue(controlID); ok {
			if s, ok := v.(string); ok {
				l.SetSelectedIdsNoRefresh([]string{l.ID() + "_" + s})
			}
		}
	} else {
		if v, ok := ctx.FormValue(controlID); ok {
			l.SetSelectedIdsNoRefresh([]string{v})
		}
	}
}
