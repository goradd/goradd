package control

import (
	"context"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
)

type RadioListI interface {
	control.RadioListI
}

// RadioList is a multi-select control that presents its choices as a list of radio buttons.
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
	l.Self = l
	l.Init(parent, id)
	return l
}

func (l *RadioList) Init(parent page.ControlI, id string) {
	l.RadioList.Init(parent, id)
	l.SetLabelDrawingMode(html5tag.LabelAfter)
	l.SetRowClass("row")
}

func (l *RadioList) this() RadioListI {
	return l.Self.(RadioListI)
}

func (l *RadioList) SetIsInline(i bool) {
	l.isInline = i
}

// SetCellClass sets a string that is applied to every cell. This is useful for setting responsive breakpoints
func (l *RadioList) SetCellClass(c string) {
	l.cellClass = c
}

// TODO: Use bootstrap styling for the columns rather than table styling
// Also coordinate with FormFieldset


// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *RadioList) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx) // skip default checkbox list attributes
	a.SetData("grctl", "bs-RadioList")
	return a
}

// RenderItem is called by the framework to render a single item in the list.
func (l *RadioList) RenderItem(item *control.ListItem) (h string) {
	selected := l.SelectedItem().ID() == item.ID()
	h = renderItemControl(item, "radio", selected, l.ID())
	h = renderCell(item, h, l.ColumnCount(), l.isInline, l.cellClass)
	return
}

func (l *RadioList) Serialize(e page.Encoder) {
	l.RadioList.Serialize(e)

	if err := e.Encode(l.isInline); err != nil {
		panic(err)
	}
	if err := e.Encode(l.cellClass); err != nil {
		panic(err)
	}
}


func (l *RadioList) Deserialize(d page.Decoder) {
	l.RadioList.Deserialize(d)

	if err := d.Decode(&l.isInline); err != nil {
		panic(err)
	}
	if err := d.Decode(&l.cellClass); err != nil {
		panic(err)
	}
}


type RadioListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []control.ListValue
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
	// OnChange is the action to take when any of the radio buttons in the list change
	OnChange action.ActionI
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
	sub := control.RadioListCreator{
		ID: c.ID,
		Items: c.Items,
		DataProvider: c.DataProvider,
		ColumnCount: c.ColumnCount,
		LayoutDirection: c.LayoutDirection,
		LabelDrawingMode: c.LabelDrawingMode,
		IsScrolling: c.IsScrolling,
		RowClass: c.RowClass,
		Value: c.Value,
		OnChange: c.OnChange,
		SaveState: c.SaveState,
		ControlOptions: c.ControlOptions,

	}
	sub.Init(ctx, ctrl)
}

// GetRadioList is a convenience method to return the control with the given id from the page.
func GetRadioList(c page.ControlI, id string) *RadioList {
	return c.Page().GetControl(id).(*RadioList)
}

func init() {
	page.RegisterControl(&RadioList{})
}