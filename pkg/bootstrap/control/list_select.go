package control

import (
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/bootstrap/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

type SelectListI interface {
	control.SelectListI
}
type SelectList struct {
	control.SelectList
}

func NewSelectList(parent page.ControlI, id string) *SelectList {
	t := new(SelectList)
	t.Init(t, parent, id)
	return t
}

func (l *SelectList) Init(self page.ControlI, parent page.ControlI, id string) {
	l.SelectList.Init(self, parent, id)
	config.LoadBootstrap(l.ParentForm())
}


func (t *SelectList) ΩDrawingAttributes(ctx context.Context) html.Attributes {
	a := t.SelectList.ΩDrawingAttributes(ctx)
	a.AddClass("form-control")
	return a
}

func init() {
	gob.RegisterName("bootstrap.selectlist", new(SelectList))
}

type SelectListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []control.ListValue
	// NilItem is a helper to add an item at the top of the list with a nil value. This is often
	// used to specify no selection, or a message that a selection is required.
	NilItem string
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// Size specifies how many items to show, and turns the list into a scrolling list
	Size int
	// Value is the initial value of the textbox. Often its best to load the value in a separate Load step after creating the control.
	Value string
	// SaveState saves the selected value so that it is restored if the form is returned to.
	SaveState bool
	page.ControlOptions
}

func (c SelectListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewSelectList(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c SelectListCreator) Init(ctx context.Context, ctrl SelectListI) {

	sub := control.SelectListCreator{
		ID: c.ID,
		Items: c.Items,
		NilItem: c.NilItem,
		DataProvider: c.DataProvider,
		Size: c.Size,
		Value: c.Value,
		SaveState: c.SaveState,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
}


// GetSelectList is a convenience method to return the control with the given id from the page.
func GetSelectList(c page.ControlI, id string) *SelectList {
	return c.Page().GetControl(id).(*SelectList)
}

func init() {
	page.RegisterControl(SelectList{})
}