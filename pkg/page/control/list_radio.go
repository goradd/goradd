package control

import (
	"bytes"
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"strings"
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
	SelectList
	// ColumnCount is the number of columns to force the list to display. It specifies the maximum number
	// of objects placed in each row wrapper. Keeping this at zero (the default)
	// will result in no row wrappers.
	ColumnCount int
	// NextItem controls how items are placed when there are columns.
	Placement   NextItemPlacement
	// LabelDrawingMode determines how labels are drawn. The default is to use the global setting.
	LabelDrawingMode html.LabelDrawingMode
	// IsScrolling determines if we are going to let the list scroll. You will need to limit the size of the
	// control for scrolling to happen.
	IsScrolling bool

}

// NewRadioList creates a new RadioList control.
func NewRadioList(parent page.ControlI, id string) *RadioList {
	l := &RadioList{}
	l.Init(l, parent, id)
	return l
}

// Init is called by subclasses.
func (l *RadioList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.SelectList.Init(self, parent, id)
	l.Tag = "div";
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time.
// You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "radiolist")
	a.AddClass("gr-cbl")

	if l.IsScrolling {
		a.AddClass("gr-cbl-scroller")
	}

	return a
}

// ΩDrawInnerHtml is called by the framework to draw the contents of the list.
func (l *RadioList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.ΩRenderItems(l.items)
	h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table").SetID(l.ID() + "_cbl"), h)
	buf.WriteString(h)
	return nil
}

func (l *RadioList) ΩRenderItems(items []ListItemI) string {
	var hItems []string
	for _,item := range items {
		hItems = append(hItems, l.ΩRenderItem(item))
	}
	if l.ColumnCount == 0 {
		return strings.Join(hItems, "")
	}
	b := GridLayoutBuilder{}
	return b.Items(hItems).
		ColumnCount(l.ColumnCount).
		Direction(l.Placement).
		RowClass("gr-cbl-row").
		Build()
}


// ΩRenderItem is called by the framework to render a single item in the list.
func (l *RadioList) ΩRenderItem(item ListItemI) (h string) {
	h = l.ΩRenderItemControl(item, "radio")
	h = l.ΩRenderCell(item, h)
	return
}

func (l *RadioList) ΩRenderItemControl(item ListItemI, typ string) string {
	labelMode := l.LabelDrawingMode
	if labelMode == html.LabelDefault {
		labelMode = page.DefaultCheckboxLabelDrawingMode
	}

	attributes := html.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", l.ID())
	attributes.Set("value", item.ID())
	attributes.Set("type", typ)
	if l.selectedId == item.ID() {
		attributes.Set("checked", "")
	}
	ctrl := html.RenderVoidTag("input", attributes)
	return html.RenderLabel(html.NewAttributes().Set("for", item.ID()), item.Label(), ctrl, labelMode)
}

func (l *RadioList) ΩRenderCell(item ListItemI, controlHtml string) string {
	var cellClass string
	var itemId string
	if l.ColumnCount == 0 {
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



