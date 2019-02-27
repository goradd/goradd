package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type ItemDirection int

const (
	HorizontalItemDirection ItemDirection = 0
	VerticalItemDirection                 = 1
)

type CheckboxListI interface {
	control.CheckboxListI
	renderItem(item control.ListItemI) string
}

// CheckboxList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-table table, or a single-table
// scrolling list much like a standard html select list.
type CheckboxList struct {
	control.CheckboxList
	isInline      bool
}

func NewCheckboxList(parent page.ControlI, id string) *CheckboxList {
	l := &CheckboxList{}
	l.Init(l, parent, id)
	return l
}

func (l *CheckboxList) Init(self CheckboxListI, parent page.ControlI, id string) {
	l.CheckboxList.Init(self, parent, id)
	l.SetLabelDrawingMode(html.LabelAfter)
}

func (l *CheckboxList) this() CheckboxListI {
	return l.Self.(CheckboxListI)
}

func (l *CheckboxList) SetIsInline(i bool) {
	l.isInline = i
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *CheckboxList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()	// skip default checkbox list attributes
	a.SetDataAttribute("grctl", "bs-checkboxlist")
	/*
	a.AddAttributeValue("gr-cbl")

	if l.isScrolling {
		a.AddAttributeValue("gr-cbl-scroller")
	} else {
		a.AddAttributeValue("gr-cbl-table")
	}*/
	return a
}

func (l *CheckboxList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml()
	buf.WriteString(h)
	return nil
}

func (l *CheckboxList) getItemsHtml() (h string) {
	items := l.ListItems()

	var order = make([]int, len(items), len(items))

	rowCount := (len(items) - 1) / l.ColumnCount() + 1
	if l.Direction() == VerticalItemDirection {

		for i := 0; i < len(items); i++ {
			row := i % rowCount
			col := i / rowCount

			order[i] = row * l.ColumnCount() + col
		}
	} else {
		for i := 0; i < len(items); i++ {
			order[i] = i
		}
	}

	i := 0
	for row := 0; row < rowCount; row++ {
		rowHtml := ""
		for col := 0; col < l.ColumnCount(); col++ {
			rowHtml += l.this().renderItem(items[order[i]])
			i++
		}
		if l.ColumnCount() > 1 {
			h += html.RenderTag("div", html.NewAttributes().AddClass("row"), rowHtml)
		} else {
			h += rowHtml
		}
	}
	return
}


func (l *CheckboxList) renderItem(item control.ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", item.ID())
	attributes.Set("type", "checkbox")
	if l.IsIdSelected(item.ID()) {
		attributes.Set("checked", "")
	}
	attributes.AddClass("form-check-input")
	ctrl := html.RenderVoidTag("input", attributes)

	h = html.RenderLabel(html.NewAttributes().Set("for", item.ID()).AddClass("form-check-label"), item.Label(), ctrl, html.LabelAfter)
	attributes = item.Attributes().Copy()
	attributes.AddClass("form-check")
	if l.isInline {
		attributes.AddClass("form-check-inline")
	}
	h = html.RenderTag("div", attributes, h)
	return
}

