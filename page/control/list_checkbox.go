package control

import (
	"github.com/spekary/goradd/page"
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"strings"
	"goradd/config"
)

type ItemDirection int

const (
	HorizontalItemDirection ItemDirection = 0
	VerticalItemDirection = 1
)

// CheckboxList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-column table, or a single-column
// scrolling list much like a standard html select list.
type CheckboxList struct {
	MultiselectList
	columns int
	direction ItemDirection
	labelDrawingMode html.LabelDrawingMode
	isScrolling bool
}

func NewCheckboxList(parent page.ControlI) *CheckboxList {
	l := &CheckboxList{}
	l.Init(l, parent)
	return l
}


func (l *CheckboxList) Init(self page.ControlI, parent page.ControlI) {
	l.MultiselectList.Init(self, parent)
	l.Tag = "div"
	l.columns = 1
	l.labelDrawingMode = config.DefaultCheckboxLabelDrawingMode
}

func (l *CheckboxList) SetColumns(columns int) *CheckboxList {
	if l.columns <= 0 {
		panic ("Columns must be at least 1.")
	}
	l.columns = columns
	l.Refresh()
	return l
}

func (l *CheckboxList) SetDirection(direction ItemDirection) *CheckboxList {
	l.direction = direction
	l.Refresh()
	return l
}

func (l *CheckboxList) SetLabelDrawingMode(mode html.LabelDrawingMode) *CheckboxList {
	l.labelDrawingMode = mode
	l.Refresh()
	return l
}

func (l *CheckboxList) SetIsScrolling(s bool) *CheckboxList {
	l.isScrolling = s
	l.Refresh()
	return l
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
		h = html.RenderTag("div",html.NewAttributes().SetClass("gr-cbl-table"), h) + "\n"
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

func (l *CheckboxList) verticalHtml (items []ListItemI) (h string) {
	lines := l.verticalHtmlItems(items)
	if l.columns == 1 {
		return strings.Join(lines, "\n")
	} else {
		columnHeight := len(lines) / l.columns + 1
		for col := 0; col < l.columns; col++ {
			colHtml := strings.Join(lines[col * columnHeight:(col + 1)* columnHeight], "\n")
			colHtml = html.RenderTag("div", html.NewAttributes().AddClass("gr-cbl-column"), colHtml)
			h += colHtml
		}
		return
	}
}


func (l *CheckboxList) verticalHtmlItems(items []ListItemI) (h []string) {
	for _,item := range items {
		if item.HasChildItems() {
			tag := "div"
			attributes := item.Attributes().Clone()
			attributes.AddClass("gr-cbl-heading")
			subItems := l.verticalHtmlItems(item.ListItems())
			h = append(h, html.RenderTag(tag, attributes, item.Label()))
			h = append(h, subItems...)
		} else {
			h = append(h, l.renderItem("div", item))
		}
	}
	return
}

func (l *CheckboxList) renderItem(tag string, item ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetId(item.Id())
	attributes.Set("name", item.Id())
	attributes.Set("type", "checkbox")
	if l.selectedIds[item.Id()] {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	h = html.RenderLabel(html.NewAttributes().Set("for", item.Id()), item.Label(), ctrl, l.labelDrawingMode)
	attributes = item.Attributes().Clone()
	attributes.AddClass("gr-cbl-item")
	h = html.RenderTag(tag, attributes, h)
	return
}

func (l *CheckboxList) horizontalHtml (items []ListItemI) (h string) {
	var itemNum int
	var rowHtml string

	for _,item := range items {
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
			rowHtml += l.renderItem("span", item)
			itemNum ++
			if itemNum == l.columns {
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
	controlId := l.Id()

	if v,ok := ctx.CheckableValue(controlId); ok {
		l.selectedIds = map[string]bool{}
		if a,ok := v.([]interface{}); ok {
			for _,id := range a {
				l.selectedIds[controlId + "_" + id.(string)] = true
			}
		}
	}
}
