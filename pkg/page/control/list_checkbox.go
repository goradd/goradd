package control

import (
	"context"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"io"
	"strings"
)

type CheckboxListI interface {
	MultiselectListI
	SetColumnCount(int) CheckboxListI
	SetLayoutDirection(direction LayoutDirection) CheckboxListI
	SetLabelDrawingMode(mode html.LabelDrawingMode) CheckboxListI
	SetIsScrolling(s bool) CheckboxListI
	SetRowClass(c string) CheckboxListI

	RenderItems(items []*ListItem) string
	RenderItem(item *ListItem) string
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
	// rowClass is the class assigned to the div wrapper around each row.
	rowClass string
}

// NewCheckboxList creates a new CheckboxList
func NewCheckboxList(parent page.ControlI, id string) *CheckboxList {
	l := &CheckboxList{}
	l.Self = l
	l.Init(parent, id)
	return l
}

// Init is called by subclasses
func (l *CheckboxList) Init(parent page.ControlI, id string) {
	l.MultiselectList.Init(parent, id)
	l.Tag = "div"
	l.rowClass = "gr-cbl-row"
	l.labelDrawingMode = page.DefaultCheckboxLabelDrawingMode
}

func (l *CheckboxList) this() CheckboxListI {
	return l.Self.(CheckboxListI)
}

// SetColumnCount sets the number of columns to use to display the list. Items will be evenly distributed
// across the columns.
func (l *CheckboxList) SetColumnCount(columns int) CheckboxListI {
	if l.columnCount < 0 {
		panic("Columns must be at least 0.")
	}
	l.columnCount = columns
	l.Refresh()
	return l.this()
}

// ColumnCount returns the current column count.
func (l *CheckboxList) ColumnCount() int {
	return l.columnCount
}

// SetLayoutDirection specifies how items are distributed across the columns.
func (l *CheckboxList) SetLayoutDirection(direction LayoutDirection) CheckboxListI {
	l.direction = direction
	l.Refresh()
	return l.this()
}

// LayoutDirection returns the direction of how items are spread across the columns.
func (l *CheckboxList) LayoutDirection() LayoutDirection {
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

// SetRowClass sets the class to the div wrapper around each row. If blank, will be given
// a default.
func (l *CheckboxList) SetRowClass(c string) CheckboxListI {
	l.rowClass = c
	l.Refresh()
	return l.this()
}


// DrawingAttributes retrieves the tag's attributes at draw time.
// You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *CheckboxList) DrawingAttributes(ctx context.Context) html.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "checkboxlist")
	a.AddClass("gr-cbl")

	if l.isScrolling {
		a.AddClass("gr-cbl-scroller")
	}
	return a
}

// DrawInnerHtml is called by the framework to draw the contents of the list.
func (l *CheckboxList) DrawInnerHtml(ctx context.Context, w io.Writer) (err error) {
	h := l.this().RenderItems(l.items)
	h = html.RenderTag("div", html.NewAttributes().SetClass("gr-cbl-table").SetID(l.ID()+"_cbl"), h)
	_,err = io.WriteString(w, h)
	return
}

func (l *CheckboxList) RenderItems(items []*ListItem) string {
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
func (l *CheckboxList) RenderItem(item *ListItem) (h string) {
	_, selected := l.selectedValues[item.Value()]
	h = renderCheckItemControl(item, "checkbox", l.labelDrawingMode, selected, l.ID())
	h = renderCell(item, h, l.columnCount > 0)
	return
}

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (l *CheckboxList) UpdateFormValues(ctx context.Context) {
	controlID := l.ID()
	grctx := page.GetContext(ctx)

	if grctx.RequestMode() == page.Server {
		// Using name attribute to return rendered checkboxes that are turned on.
		if v, ok := grctx.FormValues(controlID); ok {
			l.SetSelectedValuesNoRefresh(v)
		}
	}
}

func (l *CheckboxList) Serialize(e page.Encoder) (err error) {
	if err = l.MultiselectList.Serialize(e); err != nil {
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

func (l *CheckboxList) Deserialize(dec page.Decoder) (err error) {
	if err = l.MultiselectList.Deserialize(dec); err != nil {
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

type CheckboxListCreator struct {
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
func (c CheckboxListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewCheckboxList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c CheckboxListCreator) Init(ctx context.Context, ctrl CheckboxListI) {
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
func GetCheckboxList(c page.ControlI, id string) *CheckboxList {
	return c.Page().GetControl(id).(*CheckboxList)
}

func init() {
	page.RegisterControl(&CheckboxList{})
}