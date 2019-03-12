package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"strings"
)

type RadioListI interface {
	SelectListI
	ΩRenderItems(items []ListItemI) string
}

// RadioList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-table table, or a single-table
// scrolling list much like a standard html select list.
type RadioList struct {
	SelectList
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

// NewRadioList creates a new RadioList control.
func NewRadioList(parent page.ControlI, id string) *RadioList {
	l := &RadioList{}
	l.Init(l, parent, id)
	return l
}

// Init is called by subclasses.
func (l *RadioList) Init(self RadioListI, parent page.ControlI, id string) {
	l.SelectList.Init(self, parent, id)
	l.Tag = "div";
	l.labelDrawingMode = page.DefaultCheckboxLabelDrawingMode
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

// SetColumnCount sets the number of columns to use to display the list. Items will be evenly distributed
// across the columns.
func (l *RadioList) SetColumnCount(columns int) *RadioList {
	if l.columnCount < 0 {
		panic("Columns must be at least 0.")
	}
	l.columnCount = columns
	l.Refresh()
	return l
}

// ColumnCount returns the current column count.
func (l *RadioList) ColumnCount() int {
	return l.columnCount
}

// SetDirection specifies how items are distributed across the columns.
func (l *RadioList) SetDirection(direction LayoutDirection) *RadioList {
	l.direction = direction
	l.Refresh()
	return l
}

// Direction returns the direction of how items are spread across the columns.
func (l *RadioList) Direction() LayoutDirection {
	return l.direction
}

// SetLabelDrawingMode indicates how labels for each of the checkboxes are drawn.
func (l *RadioList) SetLabelDrawingMode(mode html.LabelDrawingMode) *RadioList {
	l.labelDrawingMode = mode
	l.Refresh()
	return l
}

// SetIsScrolling sets whether the list will scroll if it gets bigger than its bounding box.
// You will need to style the bounding box to give it limits, or else it will simply grow as
// big as the list.
func (l *RadioList) SetIsScrolling(s bool) *RadioList {
	l.isScrolling = s
	l.Refresh()
	return l
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time.
// You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "radiolist")
	a.AddClass("gr-cbl")

	if l.isScrolling {
		a.AddClass("gr-cbl-scroller")
	}

	return a
}

// ΩDrawInnerHtml is called by the framework to draw the contents of the list.
func (l *RadioList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.this().ΩRenderItems(l.items)
	h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table").SetID(l.ID() + "_cbl"), h)
	buf.WriteString(h)
	return nil
}

func (l *RadioList) ΩRenderItems(items []ListItemI) string {
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
func (l *RadioList) ΩRenderItem(item ListItemI) (h string) {
	h = renderItemControl(item, "radio", l.labelDrawingMode, item.ID() == l.selectedId, l.ID())
	h = renderCell(item, h, l.columnCount > 0)
	return
}

func renderItemControl(item ListItemI, typ string, labelMode html.LabelDrawingMode, selected bool, name string) string {
	if labelMode == html.LabelDefault {
		labelMode = page.DefaultCheckboxLabelDrawingMode
	}

	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", name)
	attributes.Set("value", item.ID())
	attributes.Set("type", typ)
	if selected {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	return html.RenderLabel(html.NewAttributes().Set("for", item.ID()), item.Label(), ctrl, labelMode)
}

func renderCell(item ListItemI, controlHtml string, hasColumns bool) string {
	var cellClass string
	var itemId string
	if !hasColumns {
		cellClass = "gr-cbl-item"
		itemId = item.ID() + "_item"

	} else {
		cellClass = "gr-cbl-cell"
		itemId = item.ID() + "_cell"
	}
	attributes := item.Attributes().Copy()
	attributes.SetID(itemId)
	attributes.AddClass(cellClass)
	return html.RenderTag("div", attributes, controlHtml)
}

// ΩUpdateFormValues is called by the framework to tell the control to update its internal values
// based on the form values sent by the browser.
func (l *RadioList) ΩUpdateFormValues(ctx *page.Context) {
	controlID := l.ID()

	if v, ok := ctx.FormValue(controlID); ok {
		l.selectedId = v
	}
}


