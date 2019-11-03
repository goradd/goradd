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
	SetColumnCount(int) RadioListI
	SetLayoutDirection(direction LayoutDirection) RadioListI
	SetLabelDrawingMode(mode html.LabelDrawingMode) RadioListI
	SetIsScrolling(s bool) RadioListI
	SetRowClass(c string) RadioListI

	RenderItems(items []*ListItem) string
	RenderItem(item *ListItem) string

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
	// rowClass is the class assigned to the div wrapper around each row.
	rowClass string
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
	l.Tag = "div"
	l.rowClass = "gr-cbl-row"
	l.labelDrawingMode = page.DefaultCheckboxLabelDrawingMode
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

// SetColumnCount sets the number of columns to use to display the list. Items will be evenly distributed
// across the columns.
func (l *RadioList) SetColumnCount(columns int) RadioListI {
	if l.columnCount < 0 {
		panic("Columns must be at least 0.")
	}
	l.columnCount = columns
	l.Refresh()
	return l.this()
}

// ColumnCount returns the current column count.
func (l *RadioList) ColumnCount() int {
	return l.columnCount
}

// SetLayoutDirection specifies how items are distributed across the columns.
func (l *RadioList) SetLayoutDirection(direction LayoutDirection) RadioListI {
	l.direction = direction
	l.Refresh()
	return l.this()
}

// LayoutDirection returns the direction of how items are spread across the columns.
func (l *RadioList) LayoutDirection() LayoutDirection {
	return l.direction
}

// SetLabelDrawingMode indicates how labels for each of the checkboxes are drawn.
func (l *RadioList) SetLabelDrawingMode(mode html.LabelDrawingMode) RadioListI {
	l.labelDrawingMode = mode
	l.Refresh()
	return l.this()
}

// SetIsScrolling sets whether the list will scroll if it gets bigger than its bounding box.
// You will need to style the bounding box to give it limits, or else it will simply grow as
// big as the list.
func (l *RadioList) SetIsScrolling(s bool) RadioListI {
	l.isScrolling = s
	l.Refresh()
	return l.this()
}

// SetRowClass sets the class to the div wrapper around each row. If blank, will be given
// a default.
func (l *RadioList) SetRowClass(c string) RadioListI {
	l.rowClass = c
	l.Refresh()
	return l.this()
}

// DrawingAttributes retrieves the tag's attributes at draw time.
// You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) DrawingAttributes(ctx context.Context) html.Attributes {
	a := l.Control.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "radiolist")
	a.AddClass("gr-cbl")

	if l.isScrolling {
		a.AddClass("gr-cbl-scroller")
	}

	return a
}

// DrawInnerHtml is called by the framework to draw the contents of the list.
func (l *RadioList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.this().RenderItems(l.items)
	h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table").SetID(l.ID()+"_cbl"), h)
	buf.WriteString(h)
	return nil
}

func (l *RadioList) RenderItems(items []*ListItem) string {
	var hItems []string
	for _,item := range items {
		hItems = append(hItems, l.this().RenderItem(item))
	}
	if l.columnCount == 0 {
		return strings.Join(hItems, "")
	}
	b := GridLayoutBuilder{}
	return b.Items(hItems).
		ColumnCount(l.columnCount).
		Direction(l.direction).
		RowClass(l.rowClass).
		Build()
}

// RenderItem is called by the framework to render a single item in the list.
func (l *RadioList) RenderItem(item *ListItem) (h string) {
	h = renderItemControl(item, "radio", l.labelDrawingMode, item.ID() == l.selectedId, l.ID())
	h = renderCell(item, h, l.columnCount > 0)
	return
}

func renderItemControl(item *ListItem, typ string, labelMode html.LabelDrawingMode, selected bool, name string) string {
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

func renderCell(item *ListItem, controlHtml string, hasColumns bool) string {
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

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (l *RadioList) UpdateFormValues(ctx *page.Context) {
	controlID := l.ID()

	if v, ok := ctx.FormValue(controlID); ok {
		l.selectedId = v
	}
}

func (l *RadioList) Serialize(e page.Encoder) (err error) {
	if err = l.SelectList.Serialize(e); err != nil {
		return
	}
	if err = e.Encode(l.columnCount); err != nil {
		return
	}
	if err = e.Encode(l.direction); err != nil {
		return
	}
	if err = e.Encode(l.labelDrawingMode); err != nil {
		return
	}
	if err = e.Encode(l.isScrolling); err != nil {
		return
	}
	if err = e.Encode(l.rowClass); err != nil {
		return
	}
	return
}

func (l *RadioList) Deserialize(dec page.Decoder) (err error) {
	if err = l.SelectList.Deserialize(dec); err != nil {
		return
	}
	if err = dec.Decode(&l.columnCount); err != nil {
		return
	}
	if err = dec.Decode(&l.direction); err != nil {
		return
	}
	if err = dec.Decode(&l.labelDrawingMode); err != nil {
		return
	}
	if err = dec.Decode(&l.isScrolling); err != nil {
		return
	}
	if err = dec.Decode(&l.rowClass); err != nil {
		return
	}
	return
}


type RadioListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// ColumnCount specifies how many columns to show
	ColumnCount int
	// LayoutDirection determines how the items are arranged in the columns
	LayoutDirection LayoutDirection
	// LabelDrawingMode specifies how the labels on the radio buttons will be associated with the buttons
	LabelDrawingMode html.LabelDrawingMode
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
func (c RadioListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewRadioList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c RadioListCreator) Init(ctx context.Context, ctrl RadioListI) {
	if c.Items != nil {
		ctrl.AddListItems(c.Items)
	}
	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(DataBinder)
		ctrl.SetDataProvider(provider)
	}
	if c.ColumnCount != 0 {
		ctrl.SetColumnCount(c.ColumnCount)
	}
	ctrl.SetLayoutDirection(c.LayoutDirection)
	if c.LabelDrawingMode != html.LabelDefault {
		ctrl.SetLabelDrawingMode(c.LabelDrawingMode)
	}
	if c.IsScrolling {
		ctrl.SetIsScrolling(true)
	}
	if c.RowClass != "" {
		ctrl.SetRowClass(c.RowClass)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
}

// GetRadioList is a convenience method to return the control with the given id from the page.
func GetRadioList(c page.ControlI, id string) *RadioList {
	return c.Page().GetControl(id).(*RadioList)
}

func init() {
	page.RegisterControl(RadioList{})
}