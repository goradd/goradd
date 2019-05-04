package control

import (
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type RadioListI interface {
	control.RadioListI
}

// RadioList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-table table, or a single-table
// scrolling list much like a standard html select list.
type RadioList struct {
	control.RadioList
	isInline  bool
	cellClass string
}

func NewRadioList(parent page.ControlI, id string) *RadioList {
	l := &RadioList{}
	l.Init(l, parent, id)
	return l
}

func (l *RadioList) Init(self RadioListI, parent page.ControlI, id string) {
	l.RadioList.Init(self, parent, id)
	l.SetLabelDrawingMode(html.LabelAfter)
	l.SetRowClass("row")
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

func (l *RadioList) SetIsInline(i bool) {
	l.isInline = i
}

// SetColumnClass sets a string that is applied to every cell. This is useful for setting responsive breakpoints
func (l *RadioList) SetCellClass(c string) {
	l.cellClass = c
}


// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()	// skip default checkbox list attributes
	a.SetDataAttribute("grctl", "bs-RadioList")
	return a
}

// ΩRenderItem is called by the framework to render a single item in the list.
func (l *RadioList) ΩRenderItem(item control.ListItemI) (h string) {
	selected := l.SelectedItem().ID() != item.ID()
	h = renderItemControl(item, "radio", selected, l.ID())
	h = renderCell(item, h, l.ColumnCount(), l.isInline, l.cellClass)
	return
}
