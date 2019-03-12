package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"strings"
)

type CheckboxListI interface {
	MultiselectListI
	ΩRenderItems(items []ListItemI) string
}

// CheckboxList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-column table, or a single-column
// scrolling list much like a standard html select list.
type CheckboxList struct {
	MultiselectList
	// columnCount is the number of columns to force the list to display. It specifies the maximum number
	// of objects placed in each row wrapper. Keeping this at zero (the default)
	// will result in no row wrappers.
	columnCount int
	// direction controls how items are placed when there are columns.
	direction LayoutDirection
	// labelDrawingMode determines how labels are drawn. The default is to use the global setting.
	labelDrawingMode html.LabelDrawingMode
	// isScrolling determines if we are going to let the list scroll. You will need to limit the size of the
	// control for scrolling to happen.
	isScrolling bool
}

// NewCheckboxList creates a new CheckboxList
func NewCheckboxList(parent page.ControlI, id string) *CheckboxList {
	l := &CheckboxList{}
	l.Init(l, parent, id)
	return l
}

// Init is called by subclasses
func (l *CheckboxList) Init(self CheckboxListI, parent page.ControlI, id string) {
	l.MultiselectList.Init(self, parent, id)
	l.Tag = "div"
	l.labelDrawingMode = page.DefaultCheckboxLabelDrawingMode
}

func (l *CheckboxList) this() CheckboxListI {
	return l.Self.(CheckboxListI)
}

// SetColumnCount sets the number of columns to use to display the list. Items will be evenly distributed
// across the columns.
func (l *CheckboxList) SetColumnCount(columns int) *CheckboxList {
	if l.columnCount < 0 {
		panic("Columns must be at least 0.")
	}
	l.columnCount = columns
	l.Refresh()
	return l
}

// ColumnCount returns the current column count.
func (l *CheckboxList) ColumnCount() int {
	return l.columnCount
}

// SetDirection specifies how items are distributed across the columns.
func (l *CheckboxList) SetDirection(direction LayoutDirection) *CheckboxList {
	l.direction = direction
	l.Refresh()
	return l
}

// Direction returns the direction of how items are spread across the columns.
func (l *CheckboxList) Direction() LayoutDirection {
	return l.direction
}

// SetLabelDrawingMode indicates how labels for each of the checkboxes are drawn.
func (l *CheckboxList) SetLabelDrawingMode(mode html.LabelDrawingMode) *CheckboxList {
	l.labelDrawingMode = mode
	l.Refresh()
	return l
}

// SetIsScrolling sets whether the list will scroll if it gets bigger than its bounding box.
// You will need to style the bounding box to give it limits, or else it will simply grow as
// big as the list.
func (l *CheckboxList) SetIsScrolling(s bool) *CheckboxList {
	l.isScrolling = s
	l.Refresh()
	return l
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
	}
	return a
}

// ΩDrawInnerHtml is called by the framework to draw the contents of the list.
func (l *CheckboxList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.this().ΩRenderItems(l.items)
	h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table").SetID(l.ID() + "_cbl"), h)
	buf.WriteString(h)
	return nil
}


func (l *CheckboxList) ΩRenderItems(items []ListItemI) string {
	var hItems []string
	for _,item := range items {
		hItems = append(hItems, l.ΩRenderItem(item))
	}
	if l.columnCount == 0 {
		return strings.Join(hItems, "")
	}
	b := GridLayoutBuilder{}
	return b.Items(hItems).
		ColumnCount(l.columnCount).
		Direction(l.direction).
		RowClass("gr-cbl-row").
		Build()
}


// ΩRenderItem is called by the framework to render a single item in the list.
func (l *CheckboxList) ΩRenderItem(item ListItemI) (h string) {
	_,selected := l.selectedIds[item.ID()]
	h = renderItemControl(item, "checkbox", l.labelDrawingMode, selected, l.ID())
	h = renderCell(item, h, l.columnCount > 0)
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
