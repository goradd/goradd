package control

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/control/list"
	"github.com/goradd/html5tag"
)

type ItemDirection int

const (
	HorizontalItemDirection ItemDirection = 0
	VerticalItemDirection                 = 1
)

type CheckboxListI interface {
	list.CheckboxListI
}

// CheckboxList is a multi-select control that presents its choices as a list of checkboxes.
// Styling is provided by divs and spans that you can provide css for in your style sheets. The
// goradd.css file has default styling to handle the basics. It wraps the whole thing in a div that can be set
// to scroll as well, so that the final structure can be styled like a multi-table table, or a single-table
// scrolling list much like a standard html select list.
type CheckboxList struct {
	list.CheckboxList
	isInline  bool
	cellClass string
}

func NewCheckboxList(parent page.ControlI, id string) *CheckboxList {
	l := &CheckboxList{}
	l.Self = l
	l.Init(parent, id)
	return l
}

func (l *CheckboxList) Init(parent page.ControlI, id string) {
	l.CheckboxList.Init(parent, id)
	l.SetLabelDrawingMode(html5tag.LabelAfter)
	l.SetRowClass("row")
}

func (l *CheckboxList) this() CheckboxListI {
	return l.Self.(CheckboxListI)
}

func (l *CheckboxList) SetIsInline(i bool) {
	l.isInline = i
}

// SetCellClass sets a string that is applied to every cell. This is useful for setting responsive breakpoints
func (l *CheckboxList) SetCellClass(c string) {
	l.cellClass = c
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *CheckboxList) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx) // skip default checkbox list attributes
	a.SetData("grctl", "bs-checkboxlist")
	return a
}

// RenderItem is called by the framework to render a single item in the list.
func (l *CheckboxList) RenderItem(item *list.Item) (h string) {
	selected := l.IsValueSelected(item.Value())
	h = renderItemControl(item, "checkbox", selected, l.ID())
	h = renderCell(item, h, l.ColumnCount(), l.isInline, l.cellClass)
	return
}

func renderItemControl(item *list.Item, typ string, selected bool, name string) string {
	attributes := html5tag.NewAttributes()
	attributes.SetID(item.ID())
	attributes.Set("name", name)
	attributes.Set("value", item.Value())
	attributes.Set("type", typ)
	attributes.AddClass("form-check-input")
	if selected {
		attributes.Set("checked", "")
	}
	ctrl := html5tag.RenderVoidTag("input", attributes)
	return html5tag.RenderLabel(html5tag.NewAttributes().Set("for", item.ID()).AddClass("form-check-label"), item.Label(), ctrl, html5tag.LabelAfter)
}

func renderCell(item *list.Item, controlHtml string, columnCount int, isInline bool, cellClass string) string {
	attributes := item.Attributes().Copy()
	attributes.SetID(item.ID() + "_item")
	attributes.AddClass("form-check")
	if isInline {
		attributes.AddClass("form-check-inline")
	}
	if columnCount > 0 {
		attributes.AddClass(fmt.Sprintf("col-%d", 12/columnCount))
	}
	if cellClass != "" {
		attributes.AddClass(cellClass)
	}
	return html5tag.RenderTag("div", attributes, controlHtml)
}

func (l *CheckboxList) Serialize(e page.Encoder) {
	l.CheckboxList.Serialize(e)

	if err := e.Encode(l.isInline); err != nil {
		panic(err)
	}
	if err := e.Encode(l.cellClass); err != nil {
		panic(err)
	}

	return
}

func (l *CheckboxList) Deserialize(d page.Decoder) {
	l.CheckboxList.Deserialize(d)

	if err := d.Decode(&l.isInline); err != nil {
		panic(err)
	}
	if err := d.Decode(&l.cellClass); err != nil {
		panic(err)
	}

	return
}

type CheckboxListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []list.ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// ColumnCount specifies how many columns to show
	ColumnCount int
	// LayoutDirection determines how the items are arranged in the columns
	LayoutDirection control.LayoutDirection
	// LabelDrawingMode specifies how the labels on the radio buttons will be associated with the buttons
	LabelDrawingMode html5tag.LabelDrawingMode
	// IsScrolling will give the inner div a vertical scroll style. You will need to style the height of the outer control to have a fixed style as well.
	IsScrolling bool
	// RowClass is the class assigned to each row
	RowClass string
	// Value is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Value string
	// SaveState saves the selected value so that it is restored if the form is returned to.
	SaveState bool
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c CheckboxListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewCheckboxList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c CheckboxListCreator) Init(ctx context.Context, ctrl CheckboxListI) {
	sub := list.CheckboxListCreator{
		ID:               c.ID,
		Items:            c.Items,
		DataProvider:     c.DataProvider,
		ColumnCount:      c.ColumnCount,
		LayoutDirection:  c.LayoutDirection,
		LabelDrawingMode: c.LabelDrawingMode,
		IsScrolling:      c.IsScrolling,
		RowClass:         c.RowClass,
		Value:            c.Value,
		SaveState:        c.SaveState,
		ControlOptions:   c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
}

// GetCheckboxList is a convenience method to return the control with the given id from the page.
func GetCheckboxList(c page.ControlI, id string) *CheckboxList {
	return c.Page().GetControl(id).(*CheckboxList)
}

func init() {
	page.RegisterControl(&CheckboxList{})
}
