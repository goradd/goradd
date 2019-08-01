package control

import (
	"bytes"
	"context"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control/data"
	"reflect"
	"strconv"
	"strings"
)

type MultiselectListI interface {
	page.ControlI
	ItemListI
	data.DataManagerEmbedder
}

// MultiselectList is a generic list box which allows multiple selections. It is here for completeness, but is not used
// very often since it doesn't present an intuitive interface and is very browser dependent on what is presented.
// A CheckboxList is better.
type MultiselectList struct {
	page.Control
	ItemList
	data.DataManager
	selectedIds map[string]bool
}

func NewMultiselectList(parent page.ControlI, id string) *MultiselectList {
	l := &MultiselectList{}
	l.Init(l, parent, id)
	return l
}

func (l *MultiselectList) Init(self MultiselectListI, parent page.ControlI, id string) {
	l.Control.Init(self, parent, id)
	l.ItemList = NewItemList(l)
	l.selectedIds = map[string]bool{}
	l.Tag = "select"
}

func (l *MultiselectList) this() MultiselectListI {
	return l.Self.(MultiselectListI)
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

func (l *MultiselectList) Validate(ctx context.Context) bool {
	if l.IsRequired() && len(l.selectedIds) == 0 {
		if l.ErrorForRequired == "" {
			l.SetValidationError(l.ΩT("A selection is required"))
		} else {
			l.SetValidationError(l.ErrorForRequired)
		}
		return false
	}
	return true
}

// ΩUpdateFormValues is an internal function that lets us reflect the value of the selection on the web override
func (l *MultiselectList) ΩUpdateFormValues(ctx *page.Context) {
	id := l.ID()

	if a, ok := ctx.FormValues(id); ok {
		l.selectedIds = map[string]bool{}
		for _, v := range a {
			l.selectedIds[v] = true
		}
	}
}

func (l *MultiselectList) SelectedItems() []ListItemI {
	items := []ListItemI{}
	if len(l.selectedIds) == 0 {
		return nil
	}
	for id := range l.selectedIds {
		item := l.GetItem(id)
		if item != nil {
			items = append(items, item)
		}
	}
	return items
}

// SetSelectedIds sets the current selection to the given ids. You must ensure that the items with the ids exist, it will
// not attempt to make sure the items exist.
func (l *MultiselectList) SetSelectedIds(ids []string) {
	l.SetSelectedIdsNoRefresh(ids)
	l.Refresh()
}

func (l *MultiselectList) SetSelectedIdsNoRefresh(ids []string) {
	l.selectedIds = map[string]bool{}

	if ids == nil {
		return
	}

	for _, id := range ids {
		l.selectedIds[id] = true
	}
}

func (l *MultiselectList) SetSelectedIdNoRefresh(id string, value bool) {
	if value {
		l.selectedIds[id] = true
	} else {
		delete(l.selectedIds, id)
	}
}

// Value implements the Valuer interface for general purpose value getting and setting
func (l *MultiselectList) Value() interface{} {
	return l.SelectedIds()
}

// SetValue implements the Valuer interface for general purpose value getting and setting
func (l *MultiselectList) SetValue(v interface{}) {
	l.selectedIds = map[string]bool{}
	switch ids := v.(type) {
	case string:
		a := strings.Split(ids, ",")
		for _, v := range a {
			l.selectedIds[v] = true
		}

	case []string:
		for _, v := range ids {
			l.selectedIds[v] = true
		}

	case ListItemI:
		l.selectedIds[ids.ID()] = true

	case []ListItemI:
		for _, v := range ids {
			l.selectedIds[v.ID()] = true
		}

	default:
		panic("Unknown id list type")
	}
}

// SelectedIds returns a list of ids sorted by id number that correspond to the selection
func (l *MultiselectList) SelectedIds() []string {
	ids := make([]string, 0, len(l.selectedIds))
	for id := range l.selectedIds {
		ids = append(ids, id)
	}
	SortIds(ids)
	return ids
}

func (l *MultiselectList) SelectedLabels() []string {
	labels := []string{}

	for _, id := range l.SelectedIds() {
		item := l.GetItem(id)
		if item != nil {
			labels = append(labels, item.Label())
		}
	}
	return labels
}

func (l *MultiselectList) SelectedValues() []interface{} {
	values := []interface{}{}

	for _, id := range l.SelectedIds() {
		item := l.GetItem(id)
		if item != nil {
			values = append(values, item.Value())
		}
	}
	return values
}

// SetData overrides the default data setter to add objects to the item list.
// The result is kept in memory currently.
// ItemLister, ItemIDer, Labeler or Stringer types. This function can accept one or more lists of items, or
// single items.
func (l *MultiselectList) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic("you must call SetData with a slice or array")
	}

	l.ItemList.Clear()
	l.AddListItems(data)
}

// ΩMarshalState is an internal function to save the state of the control
func (l *MultiselectList) ΩMarshalState(m maps.Setter) {
	var ids = []string{}
	for id := range l.selectedIds {
		ids = append(ids, id)
	}
	m.Set("sel", ids)
}

// ΩUnmarshalState is an internal function to restore the state of the control
func (l *MultiselectList) ΩUnmarshalState(m maps.Loader) {
	l.selectedIds = map[string]bool{}

	if s, ok := m.Load("sel"); ok {
		if ids, ok := s.([]string); ok {
			for _, id := range ids {
				l.selectedIds[id] = true
			}
		}
	}
}

// ΩDrawingAttributes retrieves the tag's attributes at draw time. You should not normally need to call this, and the
// attributes are disposed of after drawing, so they are essentially read-only.
func (l *MultiselectList) ΩDrawingAttributes() *html.Attributes {
	a := l.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "multilist")
	a.Set("name", l.ID()) // needed for posts
	a.Set("multiple", "")
	if l.IsRequired() {
		a.Set("required", "")
	}
	return a
}

func (l *MultiselectList) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	if l.HasDataProvider() {
		l.LoadData(ctx, l)
	}
	h := l.getItemsHtml(l.items)
	buf.WriteString(h)
	return nil
}

func (l *MultiselectList) getItemsHtml(items []ListItemI) string {
	var h = ""

	for _, item := range items {
		if item.HasChildItems() {
			tag := "optgroup"
			innerhtml := l.getItemsHtml(item.ListItems())
			attributes := item.Attributes().Copy()
			attributes.Set("label", item.Label())
			h += html.RenderTag(tag, attributes, innerhtml) + "\n"
		} else {
			attributes := item.Attributes().Copy()
			attributes.Set("value", item.ID())
			if l.IsIdSelected(item.ID()) {
				attributes.Set("selected", "")
			}
			h += html.RenderTag("option", attributes, item.Label()) + "\n"
		}
	}
	return h
}

func (l *MultiselectList) IsIdSelected(id string) bool {
	v, ok := l.selectedIds[id]
	return ok && v
}

type MultiselectListCreator struct {
	ID string
	// Items is a static list of labels and values that will be in the list. Or, use a DataProvider to dynamically generate the items.
	Items []ListValue
	// DataProvider is the id of a control that will dynamically provide the data for the list and that implements the DataProvider interface.
	// Often this is the parent of the control.
	DataProvider string
	// Size specifies how many items to show, and turns the list into a scrolling list
	Size int
	// SaveState saves the selected value so that it is restored if the form is returned to.
	SaveState bool
	page.ControlOptions
}

func (c MultiselectListCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewMultiselectList(parent, c.ID)

	if c.Items != nil {
		ctrl.AddListItems(c.Items)
	}

	if c.DataProvider != "" {
		// If this fails, then perhaps you are giving a data provider id for a control that is not yet created. Create the control first.
		provider := parent.Page().GetControl(c.DataProvider)
		ctrl.SetDataProvider(provider.(data.DataBinder))
	}

	if c.Size != 0 {
		ctrl.SetSize(c.Size)
	}
	ctrl.ApplyOptions(c.ControlOptions)
	if c.SaveState {
		ctrl.SaveState(ctx, c.SaveState)
	}
	return ctrl
}


// GetSelectList is a convenience method to return the control with the given id from the page.
func GetMultiselectList(c page.ControlI, id string) *MultiselectList {
	return c.Page().GetControl(id).(*MultiselectList)
}
