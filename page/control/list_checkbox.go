package control

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"strings"
)

type ItemDirection int

const (
	HorizontalItemDirection ItemDirection = 0
	VerticalItemDirection                 = 1
)

type CheckboxListI interface {
	MultiselectListI
	RenderItem(tag string, item ListItemI) (h string)
}

// CheckboxList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-table table, or a single-table
// scrolling list much like a standard html select list.
type CheckboxList struct {
	MultiselectList
	columnCount      int
	direction        ItemDirection
	labelDrawingMode html.LabelDrawingMode
	isScrolling      bool
}

func NewCheckboxList(parent page.ControlI) *CheckboxList {
	l := &CheckboxList{}
	l.Init(l, parent)
	return l
}

func (l *CheckboxList) Init(self page.ControlI, parent page.ControlI) {
	l.MultiselectList.Init(self, parent)
	l.Tag = "div"
	l.columnCount = 1
	l.labelDrawingMode = page.DefaultCheckboxLabelDrawingMode
}

func (l *CheckboxList) this() CheckboxListI {
	return l.Self.(CheckboxListI)
}

func (l *CheckboxList) SetColumnCount(columns int) CheckboxListI {
	if l.columnCount <= 0 {
		panic("Columns must be at least 1.")
	}
	l.columnCount = columns
	l.Refresh()
	return l.this()
}

func (l *CheckboxList) ColumnCount() int {
	return l.columnCount
}

func (l *CheckboxList) SetDirection(direction ItemDirection) CheckboxListI {
	l.direction = direction
	l.Refresh()
	return l.this()
}

func (l *CheckboxList) Direction() ItemDirection {
	return l.direction
}


func (l *CheckboxList) SetLabelDrawingMode(mode html.LabelDrawingMode) CheckboxListI {
	l.labelDrawingMode = mode
	l.Refresh()
	return l.this()
}

func (l *CheckboxList) SetIsScrolling(s bool) CheckboxListI {
	l.isScrolling = s
	l.Refresh()
	return l.this()
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *CheckboxList) DrawingAttributes() *html.Attributes {
	a := l.Control.DrawingAttributes()
	a.SetDataAttribute("grctl", "checkboxlist")
	a.AddClass("gr-cbl")

	if l.isScrolling {
		a.AddClass("gr-cbl-scroller")
	} else {
		a.AddClass("gr-cbl-table")
	}
	return a
}

func (l *CheckboxList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	if l.isScrolling {
		h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table"), h)
	}
	buf.WriteString(h)
	return nil
}

func (l *CheckboxList) getItemsHtml(items []ListItemI) string {
	if l.direction == VerticalItemDirection {
		return l.verticalHtml(items)
	} else {
		return l.horizontalHtml(items)
	}
}

func (l *CheckboxList) verticalHtml(items []ListItemI) (h string) {
	lines := l.verticalHtmlItems(items)
	if l.columnCount == 1 {
		return strings.Join(lines, "\n")
	} else {
		columnHeight := len(lines)/l.columnCount + 1
		for col := 0; col < l.columnCount; col++ {
			colHtml := strings.Join(lines[col*columnHeight:(col+1)*columnHeight], "\n")
			colHtml = html.RenderTag("div", html.NewAttributes().AddClass("gr-cbl-table"), colHtml)
			h += colHtml
		}
		return
	}
}

func (l *CheckboxList) verticalHtmlItems(items []ListItemI) (h []string) {
	for _, item := range items {
		if item.HasChildItems() {
			tag := "div"
			attributes := item.Attributes().Clone()
			attributes.AddClass("gr-cbl-heading")
			subItems := l.verticalHtmlItems(item.ListItems())
			h = append(h, html.RenderTag(tag, attributes, item.Label()))
			h = append(h, subItems...)
		} else {
			h = append(h, l.this().RenderItem("div", item))
		}
	}
	return
}

func (l *CheckboxList) RenderItem(tag string, item ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", item.ID())
	attributes.Set("type", "checkbox")
	if l.IsIdSelected(item.ID()) {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	h = html.RenderLabel(html.NewAttributes().Set("for", item.ID()), item.Label(), ctrl, l.labelDrawingMode)
	attributes = item.Attributes().Clone()
	attributes.SetID(item.ID() + "_cell")
	attributes.AddClass("gr-cbl-item")
	h = html.RenderTag(tag, attributes, h)
	return
}

func (l *CheckboxList) horizontalHtml(items []ListItemI) (h string) {
	var itemNum int
	var rowHtml string

	for _, item := range items {
		if item.HasChildItems() {
			if itemNum != 0 {
				// output a row
				h += html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-row"), rowHtml)
				rowHtml = ""
				itemNum = 0
			}
			tag := "div"
			attributes := item.Attributes().Clone()
			attributes.AddClass("gr-cbl-heading")
			h += html.RenderTag(tag, attributes, item.Label())
			h += l.horizontalHtml(item.ListItems())
		} else {
			rowHtml += l.this().RenderItem("span", item)
			itemNum++
			if itemNum == l.columnCount {
				// output a row
				h += html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-row"), rowHtml)
				rowHtml = ""
				itemNum = 0
			}
		}
	}
	if itemNum != 0 {
		h += html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-row"), rowHtml)
	}
	return
}

func (l *CheckboxList) UpdateFormValues(ctx *page.Context) {
	controlID := l.ID()

	if v, ok := ctx.CheckableValue(controlID); ok {
		l.selectedIds = map[string]bool{}
		if a, ok := v.([]interface{}); ok {
			for _, id := range a {
				l.selectedIds[controlID+"_"+id.(string)] = true
			}
		}
	}
}
