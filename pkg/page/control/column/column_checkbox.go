package column

import (
	"context"
	"encoding/gob"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/page/event"
)

const (
	AllClickAction = iota + 1000
)

type CheckboxColumnI interface {
	control.ColumnI
	CheckboxAttributes(data interface{}) html.Attributes
}

// CheckboxColumn is a table column that contains a checkbox in each row.
// You must provide it a CheckboxProvider to connect ids and default data
// to the checkbox. Use Changes() to get the list of checkbox ids that have changed since the list was initially drawn.
type CheckboxColumn struct {
	control.ColumnBase
	showCheckAll bool
	checkboxer   CheckboxProvider
	current      map[string]bool // currently displayed items
	changes      map[string]bool // changes recorded
}

// NewChecboxColumn creates a new table column that contains a checkbox. You must provide
// a CheckboxProvider which will connect checkbox states to data states
//
// The table will keep track of what checkboxes have been clicked and the new values. Call Changes() to get those
// changes. Or, if you are recording your changes in real time, attach a CheckboxColumnClick event to the table.
func NewCheckboxColumn(p CheckboxProvider) *CheckboxColumn {
	if p == nil {
		panic("a checkbox attribute provider is required")
	}

	i := CheckboxColumn{checkboxer: p}
	i.Init()
	return &i
}

func (col *CheckboxColumn) Init() {
	col.ColumnBase.Init(col)
	col.SetIsHtml(true)
	col.changes = make(map[string]bool)
}

func (col *CheckboxColumn) this() CheckboxColumnI {
	return col.Self.(CheckboxColumnI)
}

// SetShowCheckAll will cause the CheckAll checkbox to appear in the header. You must show at least one header
// row to see the checkboxes too.
func (col *CheckboxColumn) SetShowCheckAll(s bool) *CheckboxColumn {
	col.showCheckAll = s
	return col
}

// HeaderCellHtml is called by the Table drawing system to draw the HeaderCellHtml.
func (col *CheckboxColumn) HeaderCellHtml(ctx context.Context, rowNum int, colNum int) (h string) {
	if col.showCheckAll {
		a := col.this().CheckboxAttributes(nil)
		a.Set("type", "checkbox")
		h += html.RenderVoidTag("input", a)
	}
	if col.IsSortable() {
		h += col.RenderSortButton(col.Title())
	} else if col.Title() != "" {
		h += col.Title()
	}

	return
}

// CheckboxAttributes returns the attributes for the input tag that will display the checkbox.
// If data is nil, it indicates a checkAll box.
func (col *CheckboxColumn) CheckboxAttributes(data interface{}) html.Attributes {
	p := col.checkboxer
	a := p.Attributes(data)
	if a == nil {
		a = html.NewAttributes()
	}
	var id string
	var pubid string

	if data == nil {
		pubid = col.ParentTable().ID() + "_" + col.ID() + "_all"
		a.Set("id", pubid)
		a.SetDataAttribute("grAll", "1")
	} else if id = p.RowID(data); id != "" {
		// TODO: optionally encrypt the id in case its a database id. Difficult since database ids might themselves be large hashes (aka Google data store)
		// Perhaps use the checkbox provider to do that?
		pubid = col.ParentTable().ID() + "_" + col.ID() + "_" + id
		a.Set("id", pubid)
		a.SetDataAttribute("grCheckcol", "1")
		col.current[id] = p.IsChecked(data)
		a.Set("name", col.ParentTable().ID()+"_"+col.ID())
		a.Set("value", id)
	} else {
		panic("A checkbox id is required.")
	}

	if newVal, ok := col.changes[id]; ok { // If we have recorded a change, use this value on refresh.
		if newVal {
			a.Set("checked", "")
		}
	} else if p.IsChecked(data) { // otherwise, use the data value
		a.Set("checked", "")
	}
	return a
}

// CellText is called by the Table drawing mechanism to draw the content of a cell, which in this case will be
// a checkbox.
func (col *CheckboxColumn) CellText(ctx context.Context, rowNum int, colNum int, data interface{}) string {
	a := col.this().CheckboxAttributes(data)
	a.Set("type", "checkbox")
	return html.RenderVoidTag("input", a)
}

// Changes returns a map of ids corresponding to checkboxes that have changed. Both true and false values indicate the
// current state of that particular checkbox. Note that if a user checks a box, then checks it again, even though it
// is back to its original value, it will still show up in the changes list.
func (col *CheckboxColumn) Changes() map[string]bool {
	return col.changes
}

// ResetChanges resets the column so it is ready to accept new data. You might need to call this if you have previously
// called SaveState. Or, change DataID in the CheckboxProvider to cause the changes to reset.
func (col *CheckboxColumn) ResetChanges() {
	col.changes = make(map[string]bool)
}

// UpdateFormValues will look for changes to our checkboxes and record those changes.
func (col *CheckboxColumn) UpdateFormValues(ctx *page.Context) {
	if ctx.RequestMode() == page.Server {
		// Using standard form submission rules. Only ON checkboxes get sent to us, so we have to figure out what got turned off
		recent := make(map[string]bool)
		if values, ok := ctx.FormValues(col.ParentTable().ID() + "_" + col.ID()); ok {
			for _, value := range values {
				recent[value] = true
			}
		}
		// otherwise its as if nothing was checked, which might happen if everything got turned off

		for k, v := range col.current {
			if _, ok := recent[k]; ok {
				// set to true
				if !v {
					col.changes[k] = true
				} else {
					// same value as original
					delete(col.changes, k)
				}
			} else {
				// set to false
				if v {
					col.changes[k] = false
				} else {
					delete(col.changes, k)
				}
			}
		}
	} else {
		// We just get notified of the ids of checkboxes that changed since the last time we checked
		for k, v := range col.current {
			if v2, ok := ctx.FormValue(col.ParentTable().ID() + "_" + col.ID() + "_" + k); ok {
				b2 := page.ConvertToBool(v2)
				if v != b2 {
					col.changes[k] = b2
				} else {
					// same value as original
					delete(col.changes, k)
				}
			}
		}
	}
}

// AddActions adds actions to the table that the column can respond to.
func (col *CheckboxColumn) AddActions(t page.ControlI) {
	t.On(event.CheckboxColumnClick().Selector(`input[data-gr-all]`), action.Ajax(col.ParentTable().ID() + "_" + col.ID(), control.ColumnAction).ActionValue(AllClickAction), action.PrivateAction{})
}

// Action is called by the framework to respond to an event. Here it responds to a click in the CheckAll box.
func (col *CheckboxColumn) Action(ctx context.Context, params page.ActionParams) {
	switch params.ActionValueInt() {
	case AllClickAction:
		p := new(event.CheckboxColumnActionValues)
		ok, err := params.EventValue(p)
		if ok && err == nil {
			col.allClick(p.Id, p.Checked, p.Row, p.Column)
		}
	}
}

// The check all checkbox has been checked.
func (col *CheckboxColumn) allClick(id string, checked bool, rowNum int, colNum int) {
	all := col.checkboxer.All()

	// if we have a checkboxer that will help us check all the objects in the table, use it
	if all != nil {
		for k, v := range all {
			if v == checked {
				col.changes[k] = checked
			} else {
				delete(col.changes, k)
			}
		}
		// Fire javascript to check all visible
		col.ParentTable().ParentForm().Response().ExecuteSelectorFunction(`input[data-gr-checkcol]`, `prop`, page.PriorityStandard, `checked`, checked)

	} else {
		// Fire javascript to check all visible and trigger a change
		if checked {
			col.ParentTable().ParentForm().Response().ExecuteSelectorFunction(`input[data-gr-checkcol]:not(:checked)`, `click`, page.PriorityStandard)

		} else {
			col.ParentTable().ParentForm().Response().ExecuteSelectorFunction(`input[data-gr-checkcol]:checked`, `click`, page.PriorityStandard)
		}
	}

}

// PreRender is called by the Table to tell the column that it is about to draw. Here we are resetting the list of
// currently showing checkboxes so that we can keep track of what is displayed. This is required to keep track of
// which boxes are checked in the event that Javascript is off.
func (col *CheckboxColumn) PreRender() {
	col.current = make(map[string]bool)
}

// MarshalState is an internal function to save the state of the control
func (t *CheckboxColumn) MarshalState(m maps.Setter) {
	m.Set(t.ID()+"_changes", t.changes)
	m.Set(t.ID()+"_dataid", t.checkboxer.DataID())
}

// UnmarshalState is an internal function to restore the state of the control
func (t *CheckboxColumn) UnmarshalState(m maps.Loader) {
	if v, ok := m.Load(t.ID() + "_dataid"); ok {
		if dataid, ok := v.(string); ok {
			if dataid == t.checkboxer.DataID() { // only restore checkboxes if the data itself has not changed
				if v, ok := m.Load(t.ID() + "_changes"); ok {
					if s, ok := v.(map[string]bool); ok {
						t.changes = s
					}
				}
			}
		}
	}
}

// The CheckboxProvider interface defines a set of functions that you implement to provide for the initial display
// of a checkbox. You can descend your own CheckboxProvider from the DefaultCheckboxProvider to get the default
// behavior, and then add whatever functions you need to implement. Be sure to register your custom provider
// with gob.
type CheckboxProvider interface {
	// RowID should return a unique id corresponding to the given data item. It is used to track the checked state of an individual checkbox.
	RowID(data interface{}) string
	// IsChecked should return true if the checkbox corresponding to the row data should initially be checked. After the
	// initial draw, the table will keep track of the state of the checkbox, meaning you do not need to live update your data.
	// If you are using the table just as a selection of items to act on, just return false here.
	IsChecked(data interface{}) bool
	// Attributes returns the attributes that will be applied to the checkbox corresponding to the data row.
	// Use this primarily for providing custom attributes. Return nil if you have no custom attributes.
	Attributes(data interface{}) html.Attributes
	// If you enable the checkAll box, you can use this to return a map of all the ids and their initial values here. This is
	// mostly helpful if your table is not showing all the rows at once (i.e. you are using a paginator or scroller and
	// only showing a subset of data at one time). If your table is showing a checkAll box, and you return nil here, the
	// checkAll will only perform a javascript checkAll, and thus only check the visible items.
	All() map[string]bool
	// DataID should return an id that identifies the overall data. This could be a database record id.
	// It is used to determine if the checkboxes in the column should be reset if SaveState is on.
	// If the DataID changes, and SaveState is on, it will reset the changes.
	DataID() string
}

// The DefaultCheckboxProvider is a mixin you can use to base your CheckboxProvider, and that will provide default
// functionality for the methods you don't want to implement.
type DefaultCheckboxProvider struct{}

func (c DefaultCheckboxProvider) DataID() string {
	return ""
}

func (c DefaultCheckboxProvider) RowID(data interface{}) string {
	return ""
}

func (c DefaultCheckboxProvider) IsChecked(data interface{}) bool {
	return false
}

func (c DefaultCheckboxProvider) Attributes(data interface{}) html.Attributes {
	return nil
}

func (c DefaultCheckboxProvider) All() map[string]bool {
	return nil
}

func init() {
	gob.Register(map[string]bool(nil)) // We must register this here because we are putting the changes map into the session,
	gob.Register(DefaultCheckboxProvider{}) // We must register this here because we are putting the changes map into the session,
}

// CheckboxColumnCreator creates a column of checkboxes.
type CheckboxColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// ShowCheckAll will show a checkbox in the header that the user can use to check all the boxes in the column.
	ShowCheckAll bool
	// CheckboxProvider tells us which checkboxes are on or off, and how the checkboxes are styled.
	CheckboxProvider   CheckboxProvider
	// Title is the title of the column that appears in the header
	Title string
	// Sortable makes the column display sort arrows in the header
	Sortable bool
	control.ColumnOptions
}

func (c CheckboxColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	col := NewCheckboxColumn(c.CheckboxProvider)
	if c.ID != "" {
		col.SetID(c.ID)
	}
	if c.ShowCheckAll {
		col.SetShowCheckAll(true)
	}
	if c.Title != "" {
		col.SetTitle(c.Title)
	}
	if c.Sortable {
		col.SetSortable()
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}