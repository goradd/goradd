package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"strings"
)

type ItemDirection int

const (
	HorizontalItemDirection ItemDirection = 0
	VerticalItemDirection                 = 1
)

type CheckboxListI interface {
	MultiselectListI
	ΩRenderItem(tag string, item ListItemI) (h string)
}

// CheckboxList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-column table, or a single-column
// scrolling list much like a standard html select list.
type CheckboxList struct {
	MultiselectList
	columnCount      int
	direction        ItemDirection
	labelDrawingMode html.LabelDrawingMode
	isScrolling      bool
}

// NewCheckboxList creates a new CheckboxList
func NewCheckboxList(parent page.ControlI, id string) *CheckboxList {
	l := &CheckboxList{}
	l.Init(l, parent, id)
	return l
}

// Init is called by subclasses
func (l *CheckboxList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.MultiselectList.Init(self, parent, id)
	l.Tag = "div"
	l.columnCount = 1
	l.labelDrawingMode = page.DefaultCheckboxLabelDrawingMode
}

func (l *CheckboxList) this() CheckboxListI {
	return l.Self.(CheckboxListI)
}

// SetColumnCount sets the number of columns to use to display the list. Items will be evenly distributed
// across the columns.
func (l *CheckboxList) SetColumnCount(columns int) CheckboxListI {
	if l.columnCount <= 0 {
		panic("Columns must be at least 1.")
	}
	l.columnCount = columns
	l.Refresh()
	return l.this()
}

// ColumnCount returns the current column count.
func (l *CheckboxList) ColumnCount() int {
	return l.columnCount
}

// SetDirection specifies how items are distributed across the columns.
func (l *CheckboxList) SetDirection(direction ItemDirection) CheckboxListI {
	l.direction = direction
	l.Refresh()
	return l.this()
}

// Direction returns the direction of how items are spread across the columns.
func (l *CheckboxList) Direction() ItemDirection {
	return l.direction
}

// SetLabelDrawingMode indicates how labels for each of the checkboxes are drawn.
func (l *CheckboxList) SetLabelDrawingMode(mode html.LabelDrawingMode) CheckboxListI {
	l.labelDrawingMode = mode
	l.Refresh()
	return l.this()
}

// SetIsScrolling sets whether the list will scroll if it gets bigger than its bounding box.
// You will need to style the bounding box to give it limits, or else it will simply grow as
// big as the list.
func (l *CheckboxList) SetIsScrolling(s bool) CheckboxListI {
	l.isScrolling = s
	l.Refresh()
	return l.this()
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time.
// You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *CheckboxList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "checkboxlist")
	a.AddClass("gr-cbl")

	if l.isScrolling {
		a.AddClass("gr-cbl-scroller")
	} else {
		a.AddClass("gr-cbl-table")
	}
	return a
}

// ΩDrawInnerHtml is called by the framework to draw the contents of the list.
func (l *CheckboxList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	if l.isScrolling {
		h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table").SetID(l.ID() + "_cbl"), h)
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
			attributes := item.Attributes().Copy()
			attributes.AddClass("gr-cbl-heading")
			subItems := l.verticalHtmlItems(item.ListItems())
			h = append(h, html.RenderTag(tag, attributes, item.Label()))
			h = append(h, subItems...)
		} else {
			h = append(h, l.this().ΩRenderItem("div", item))
		}
	}
	return
}

// ΩRenderItem draws an item in the list. You do not normally need to call this, but subclasses can override it.
func (l *CheckboxList) ΩRenderItem(tag string, item ListItemI) (h string) {
	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", l.ID())
	attributes.Set("value", item.ID())
	attributes.Set("type", "checkbox")
	if l.IsIdSelected(item.ID()) {
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
			attributes := item.Attributes().Copy()
			attributes.AddClass("gr-cbl-heading")
			h += html.RenderTag(tag, attributes, item.Label())
			h += l.horizontalHtml(item.ListItems())
		} else {
			rowHtml += l.this().ΩRenderItem("span", item)
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

// ΩUpdateFormValues is called by the framework to tell the control to update its internal values
// based on the form values sent by the browser.
func (l *CheckboxList) ΩUpdateFormValues(ctx *page.Context) {
	controlID := l.ID()

	if ctx.RequestMode() == page.Server {
		// Using name attribute to return rendered checkboxes that are turned on.
		if v, ok := ctx.FormValues(controlID); ok {
			l.SetSelectedIdsNoRefresh(v)
		}
	} else {
		// Ajax will only send changed items based on their ids
		for _,item := range l.ListItems() {
			if v, ok := ctx.FormValue(item.ID()); ok {
				l.SetSelectedIdNoRefresh(item.ID(), page.ConvertToBool(v))
			}
		}
	}
}
