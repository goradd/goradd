package list

import (
	"context"
	reflect2 "github.com/goradd/goradd/pkg/any"
	control2 "github.com/goradd/goradd/pkg/page/control"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

type MultiselectListI interface {
	page.ControlI
	ListI
	control2.DataManagerI
}

// MultiselectList is a generic list box which allows multiple selections. It is here for completeness, but is not used
// very often since it doesn't present an intuitive interface and is very browser dependent on what is presented.
// A CheckboxList is better.
type MultiselectList struct {
	page.ControlBase
	List
	control2.DataManager
	selectedValues map[string]bool
}

func NewMultiselectList(parent page.ControlI, id string) *MultiselectList {
	l := &MultiselectList{}
	l.Init(l, parent, id)
	return l
}

func (l *MultiselectList) Init(self any, parent page.ControlI, id string) {
	l.ControlBase.Init(self, parent, id)
	l.List = NewList(l)
	l.selectedValues = map[string]bool{}
	l.Tag = "select"
}

func (l *MultiselectList) this() MultiselectListI {
	return l.Self().(MultiselectListI)
}

func (l *MultiselectList) SetSize(size int) MultiselectListI {
	l.SetAttribute("size", strconv.Itoa(size))
	l.Refresh()
	return l.this()
}

func (l *MultiselectList) Size() int {
	a := l.Attribute("size")
	if a == "" {
		return 0
	} else {
		s, err := strconv.Atoi(a)
		if err != nil {
			return 0
		}
		return s
	}
}

func (l *MultiselectList) Validate(_ context.Context) bool {
	if l.IsRequired() && len(l.selectedValues) == 0 {
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
func (l *MultiselectList) UpdateFormValues(ctx context.Context) {
	id := l.ID()

	if a, ok := page.GetContext(ctx).FormValues(id); ok {
		l.selectedValues = map[string]bool{}
		for _, v := range a {
			l.selectedValues[v] = true
		}
	}
}

func (l *MultiselectList) SelectedItems() []*Item {
	var items []*Item

	if len(l.selectedValues) == 0 {
		return nil
	}
	for v := range l.selectedValues {
		_, item := l.GetItemByValue(v)
		if item != nil {
			items = append(items, item)
		}
	}
	return items
}

// SetSelectedValues sets the current selection to the given values. You must ensure that the items with the values exist, it will
// not attempt to make sure the items exist.
func (l *MultiselectList) SetSelectedValues(values []string) {
	l.SetSelectedValuesNoRefresh(values)
	l.Refresh()
}

func (l *MultiselectList) SetSelectedValuesNoRefresh(values []string) {
	l.selectedValues = map[string]bool{}

	if values == nil {
		return
	}

	for _, v := range values {
		l.selectedValues[v] = true
	}
}

func (l *MultiselectList) SetSelectedValueNoRefresh(value string, on bool) {
	if on {
		l.selectedValues[value] = true
	} else {
		delete(l.selectedValues, value)
	}
}

// Value implements the Valuer interface for general purpose value getting and setting.
func (l *MultiselectList) Value() interface{} {
	return l.SelectedValues()
}

// ValueString returns the values as a comma-delimited string.
func (l *MultiselectList) ValueString() string {
	return strings.Join(l.SelectedValues(), ",")
}

// SetValue implements the Valuer interface for general purpose value getting and setting
func (l *MultiselectList) SetValue(v interface{}) {
	l.selectedValues = map[string]bool{}
	switch values := v.(type) {
	case string:
		a := strings.Split(values, ",")
		for _, v2 := range a {
			l.selectedValues[v2] = true
		}

	case []string:
		for _, v2 := range values {
			l.selectedValues[v2] = true
		}

	case *Item:
		l.selectedValues[values.ID()] = true

	case []*Item:
		for _, v2 := range values {
			l.selectedValues[v2.ID()] = true
		}

	default:
		if v2, ok := v.(ItemIDer); ok {
			l.selectedValues[v2.ID()] = true
		} else if reflect2.IsSlice(v) {
			items := reflect2.InterfaceSlice(v)
			for _, item := range items {
				if v2, ok := item.(ItemIDer); ok {
					l.selectedValues[v2.ID()] = true
				}
			}
		} else {
			panic("Unknown id list type")
		}
	}
}

// SelectedIds returns a list of ids sorted by id number that correspond to the selection
func (l *MultiselectList) SelectedIds() []string {
	ids := make([]string, 0, len(l.selectedValues))
	for v := range l.selectedValues {
		id, _ := l.GetItemByValue(v)
		if id != "" {
			ids = append(ids, id)
		}
	}
	SortIds(ids)
	return ids
}

func (l *MultiselectList) SelectedLabels() []string {
	var labels []string

	for v := range l.selectedValues {
		_, item := l.GetItemByValue(v)
		if item != nil {
			labels = append(labels, item.Label())
		}
	}
	return labels
}

func (l *MultiselectList) SelectedValues() []string {
	var values []string

	for v := range l.selectedValues {
		values = append(values, v)
	}
	return values
}

// SetData overrides the default data setter to add objects to the item list.
// The result is kept in memory currently.
// ValueLabeler, ItemIDer, Labeler or Stringer types. This function can accept one or more lists of items, or
// single items.
func (l *MultiselectList) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic("you must call SetData with a slice or array")
	}

	l.List.Clear()
	l.AddItems(data)
}

// MarshalState is an internal function to save the state of the control
func (l *MultiselectList) MarshalState(m page.SavedState) {
	values := l.SelectedValues()
	m.Set("sel", values)
}

// UnmarshalState is an internal function to restore the state of the control
func (l *MultiselectList) UnmarshalState(m page.SavedState) {
	l.selectedValues = map[string]bool{}

	if s, ok := m.Load("sel"); ok {
		if values, ok2 := s.([]string); ok2 {
			for _, v := range values {
				l.selectedValues[v] = true
			}
		}
	}
}

func (l *MultiselectList) DrawTag(ctx context.Context, w io.Writer) {
	if l.HasDataProvider() {
		l.this().LoadData(ctx, l.this())
		defer l.ResetData()
	}
	l.ControlBase.DrawTag(ctx, w)
}

// DrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *MultiselectList) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := l.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "multilist")
	a.Set("name", l.ID()) // needed for posts
	a.Set("multiple", "")
	return a
}

func (l *MultiselectList) DrawInnerHtml(_ context.Context, w io.Writer) {
	h := l.getItemsHtml(l.items)
	page.WriteString(w, h)
	return
}

func (l *MultiselectList) getItemsHtml(items []*Item) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			tag := "optgroup"
			innerhtml := l.getItemsHtml(item.Items())
			attributes := item.Attributes().Copy()
			attributes.Set("label", item.Label())
			h += html5tag.RenderTag(tag, attributes, innerhtml) + "\n"
		} else {
			attributes := item.Attributes().Copy()
			attributes.Set("value", item.Value())
			if l.IsValueSelected(item.Value()) {
				attributes.Set("selected", "")
			}
			h += html5tag.RenderTag("option", attributes, item.Label()) + "\n"
		}
	}
	return h
}

func (l *MultiselectList) IsValueSelected(v string) bool {
	b, ok := l.selectedValues[v]
	return ok && b
}

func (l *MultiselectList) Serialize(e page.Encoder) {
	l.ControlBase.Serialize(e)
	l.List.Serialize(e)
	l.DataManager.Serialize(e)

	if err := e.Encode(l.selectedValues); err != nil {
		panic(err)
	}
}

func (l *MultiselectList) Deserialize(dec page.Decoder) {
	l.ControlBase.Deserialize(dec)
	l.List.Deserialize(dec)
	l.DataManager.Deserialize(dec)
	if err := dec.Decode(&l.selectedValues); err != nil {
		panic(err)
	}
}

type MultiselectListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider control2.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// Size specifies how many items to show, and turns the list into a scrolling list
	Size int
	// SaveState saves the selected value so that it is restored if the form is returned to.
	SaveState bool
	page.ControlOptions
}

func (c MultiselectListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewMultiselectList(parent, c.ID)

	if c.Items != nil {
		ctrl.AddItems(c.Items)
	}

	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(control2.DataBinder)
		ctrl.SetDataProvider(provider)
	}

	if c.Size != 0 {
		ctrl.SetSize(c.Size)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
	return ctrl
}

// GetMultiselectList is a convenience method to return the control with the given id from the page.
func GetMultiselectList(c page.ControlI, id string) *MultiselectList {
	return c.Page().GetControl(id).(*MultiselectList)
}

func init() {
	page.RegisterControl(&MultiselectList{})
}
