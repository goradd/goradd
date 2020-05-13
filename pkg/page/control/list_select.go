package control

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"reflect"
	"strconv"
)

type SelectListI interface {
	page.ControlI
	ItemListI
	DataManagerEmbedder
	SetValue(v interface{})
	Value() interface{}
	IntValue() int
}

// SelectList is typically a dropdown list with a single selection. Items are selected by id number, and the SelectList
// completely controls the ids in the list. Create the list by calling AddItem or AddItems to add *ListItem objects.
// Or, use the embedded DataManager to load items. Set the size attribute if you want to display it as a
// scrolling list rather than a dropdown list.
type SelectList struct {
	page.ControlBase
	ItemList
	DataManager
	selectedValue string
}

// NewSelectList creates a new select list
func NewSelectList(parent page.ControlI, id string) *SelectList {
	t := &SelectList{}
	t.Self = t
	t.Init(parent, id)
	return t
}

// Init is called by subclasses.
func (l *SelectList) Init(parent page.ControlI, id string) {
	l.ControlBase.Init(parent, id)
	l.ItemList = NewItemList(l.this())
	l.Tag = "select"
}

func (l *SelectList) this() SelectListI {
	return l.Self.(SelectListI)
}


// Validate is called by the framework to validate the contents of the control. For a SelectList,
// this is typically just checking to see if something was selected if a selection is required.
func (l *SelectList) Validate(ctx context.Context) bool {
	if v := l.ControlBase.Validate(ctx); !v {
		return false
	}

	sel := l.SelectedItem()
	if l.IsRequired() && (sel == nil ||  sel.IsEmptyValue()) {
		if l.ErrorForRequired == "" {
			l.SetValidationError(l.GT("A selection is required"))
		} else {
			l.SetValidationError(l.ErrorForRequired)
		}
		return false
	}
	return true
}

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (l *SelectList) UpdateFormValues(ctx context.Context) {
	id := l.ID()

	if v, ok := page.GetContext(ctx).FormValue(id); ok {
		l.selectedValue = v
	}
}

// SelectedItem will return the currently selected item. If no item has been selected, it will return the first item
// in the list, since that is what will be showing in the selection list, and will update its internal pointer to
// make the first item the current selection.
func (l *SelectList) SelectedItem() *ListItem {
	if l.Len() == 0 {
		return nil
	}
	if l.selectedValue == "" {
		l.selectedValue = l.items[0].Value()
		return l.items[0]
	}
	_,i := l.GetItemByValue(l.selectedValue)
	return i
}

// SetSelectedValue sets the current selection to the given id.
//
// If you are using a DataProvider, you must make sure that the value will exist in the list.
// Otherwise it will compare against the current item list and panic if the item does not exist.
func (l *SelectList) SetSelectedValue(v string) {
	if !l.HasDataProvider() {
		_,item := l.GetItemByValue(v)
		if item == nil {
			panic("Attempting to set the SelectList to a value that does not exist in the list. Value: " + v)
		}
	}
	l.selectedValue = v
	l.AddRenderScript("val", v)
}

// Value implements the Valuer interface for general purpose value getting and setting
func (l *SelectList) Value() interface{} {
	if l.selectedValue == "" {
		return nil
	} else {
		return l.selectedValue
	}
}

// SetValue implements the Valuer interface for general purpose value getting and setting
func (l *SelectList) SetValue(v interface{}) {
	l.SetSelectedValue(fmt.Sprint(v))
}

// IntValue returns the select value as an integer.
func (l *SelectList) IntValue() int {
	if l.selectedValue == "" {
		return 0
	} else {
		i,_ := strconv.Atoi(l.selectedValue)
		return i
	}
}

// StringValue returns the selected value as a string
func (l *SelectList) StringValue() string {
	return l.selectedValue
}

// SelectedLabel returns the label of the selected item
func (l *SelectList) SelectedLabel() string {
	item := l.SelectedItem()
	if item != nil {
		return item.Label()
	}
	return ""
}

// MarshalState is an internal function to save the state of the control
func (l *SelectList) MarshalState(m maps.Setter) {
	m.Set("sel", l.selectedValue)
}

// UnmarshalState is an internal function to restore the state of the control
func (l *SelectList) UnmarshalState(m maps.Loader) {
	if v, ok := m.Load("sel"); ok {
		if s, ok := v.(string); ok {
			l.selectedValue = s
		}
	}
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *SelectList) DrawingAttributes(ctx context.Context) html.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetDataAttribute("grctl", "selectlist")
	a.Set("name", l.ID()) // needed for posts
	if l.IsRequired() {
		a.Set("required", "") // required for some css frameworks, but browser validation is flaky.
							  // set the "novalidate" attribute on the form for server-side validation only.
	}
	return a
}

func (l *SelectList) DrawTag(ctx context.Context) string {
	if l.HasDataProvider() {
		l.this().LoadData(ctx, l.this())
		defer l.ResetData()
	}
	return l.ControlBase.DrawTag(ctx)
}


// DrawInnerHtml is called by the framework during drawing of the control to draw the inner html of the control
func (l *SelectList) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *SelectList) getItemsHtml(items []*ListItem) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			tag := "optgroup"
			innerhtml := l.getItemsHtml(item.ListItems())
			attributes := item.Attributes().Copy()
			attributes.Set("label", item.Label())
			h += html.RenderTag(tag, attributes, innerhtml)
		} else {
			attributes := item.Attributes().Copy()

			// TODO: add the option to encrypt values in case values are sensitive

			attributes.Set("value", item.Value())
			if l.selectedValue == item.Value() {
				attributes.Set("selected", "")
			}

			h += html.RenderTag("option", attributes, item.RenderLabel())
		}
	}
	return h
}

// SetData overrides the default data setter to add objects to the item list.
// The result is kept in memory currently.
// ItemLister, ItemIDer, Labeler or Stringer types. This function can accept one or more lists of items, or
// single items.
func (l *SelectList) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic("you must call SetData with a slice or array")
	}

	l.ItemList.Clear()
	l.AddListItems(data)
}

func (l *SelectList) Serialize(e page.Encoder) (err error) {
	if err = l.ControlBase.Serialize(e); err != nil {
		return
	}
	if err = l.ItemList.Serialize(e); err != nil {
		return
	}
	if err = l.DataManager.Serialize(e); err != nil {
		return
	}

	if err = e.Encode(l.selectedValue); err != nil {
		return
	}
	return
}

func (l *SelectList) Deserialize(dec page.Decoder) (err error) {
	if err = l.ControlBase.Deserialize(dec); err != nil {
		return
	}
	if err = l.ItemList.Deserialize(dec); err != nil {
		return
	}
	if err = l.DataManager.Deserialize(dec); err != nil {
		return
	}
	if err = dec.Decode(&l.selectedValue); err != nil {
		return
	}
	return
}


type SelectListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// NilItem is a helper to add an item at the top of the list with a nil value. This is often
	// used to specify no selection, or a message that a selection is required.
	NilItem string
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// Size specifies how many items to show, and turns the list into a scrolling list
	Size int
	// Value is the initial value of the list. Often its best to load the value in a separate Load step after creating the control.
	Value string
	// OnChange is an action to take when the user changes what is selected (as in, when the javascript change event fires).
	OnChange action.ActionI
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

	if c.NilItem != "" {
		ctrl.AddItem(c.NilItem, "")
	}

	if c.Items != nil {
		ctrl.AddListItems(c.Items)
	}

	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(DataBinder)
		ctrl.SetDataProvider(provider)
	}

	if c.Value != "" {
		ctrl.SetValue(c.Value)
	}
	if c.Size != 0 {
		ctrl.SetAttribute("size", c.Size)
	}
	if c.OnChange != nil {
		ctrl.On(event.Change(), c.OnChange)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
}


// GetSelectList is a convenience method to return the control with the given id from the page.
func GetSelectList(c page.ControlI, id string) *SelectList {
	return c.Page().GetControl(id).(*SelectList)
}

func GetSelectListI(c page.ControlI, id string) SelectListI {
	return c.Page().GetControl(id).(SelectListI)
}


func init() {
	page.RegisterControl(&SelectList{})
}