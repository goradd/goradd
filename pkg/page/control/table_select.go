package control

import (
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/html5tag"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"path"
)

// PrimaryKeyer is an interface that is often implemented by model objects.
type PrimaryKeyer interface {
	PrimaryKey() string
}

type SelectTableI interface {
	TableI
	SetSelectedID(id string) SelectTableI
	SetReselectable(r bool) SelectTableI
}

// SelectTable is a table that is row selectable. To detect a row selection, trigger on event.RowSelected
type SelectTable struct {
	Table
	selectedID string
	reselectable bool
}

func NewSelectTable(parent page.ControlI, id string) *SelectTable {
	t := &SelectTable{}
	t.Self = t
	t.Init(parent, id)
	return t
}

func (t *SelectTable) Init(parent page.ControlI, id string) {
	t.Table.Init(parent, id)
	t.ParentForm().AddJavaScriptFile(path.Join(config.AssetPrefix, "goradd", "/js/goradd-scrollIntoView.js"), false, nil)
	t.ParentForm().AddJavaScriptFile(path.Join(config.AssetPrefix, "goradd", "/js/table-select.js"), false, nil)
	t.SetAttribute("tabindex", 0); // Make the entire table focusable and selectable. This can be overridden later if needed.
	t.AddClass("gr-clickable-rows")
}

func (t *SelectTable) this() SelectTableI {
	return t.Self.(SelectTableI)
}

func (t *SelectTable) GetRowAttributes(row int, data interface{}) (a html5tag.Attributes) {
	var id string

	if t.RowStyler() != nil {
		a = t.RowStyler().TableRowAttributes(row, data)
		id = a.Get("id") // styler might be giving us an id
	} else {
		a = html5tag.NewAttributes()
	}

	// try to guess the id from the data
	if id == "" {
		switch obj := data.(type) {
		case IDer:
			id = obj.ID()
		case PrimaryKeyer:
			id = obj.PrimaryKey()
		case map[string]string:
			id, _ = obj["id"]
		case maps.StringGetter:
			id = obj.Get("id")
		}
	}
	if id != "" {
		// TODO: If configured, encrypt the id so its not publicly showing database ids
		a.SetData("id", id)
		// We need an actual id for aria features
		a.SetID(t.ID() + "_" + id)
	} else {
		a.AddClass("nosel")
	}

	a.Set("role", "option")
	return a
}

func (t *SelectTable) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := t.Table.DrawingAttributes(ctx)
	a.SetData("grctl", "selecttable")
	a.Set("role", "listbox")
	a.SetData("grWidget", "goradd.SelectTable")
	if t.selectedID != "" {
		a.SetData("grOptSelectedId", t.selectedID)
	}
	if t.reselectable {
		a.SetData("grOptReselect", "1")
	}

	return a
}

func (t *SelectTable) UpdateFormValues(ctx context.Context) {
	if data := page.GetContext(ctx).CustomControlValue(t.ID(), "selectedId"); data != nil {
		t.selectedID = fmt.Sprint(data)
	}
}

func (t *SelectTable) SelectedID() string {
	return t.selectedID
}

func (t *SelectTable) Value() interface{} {
	return t.selectedID
}

func (t *SelectTable) SetSelectedID(id string) SelectTableI {
	t.selectedID = id
	t.ExecuteWidgetFunction("option", "selectedId", id)
	return t.this()
}

// SetReselectable determines if the user can send a select command when tapping the currently selected item.
func (t *SelectTable) SetReselectable(r bool) SelectTableI {
	t.reselectable = r
	return t.this()
}


func (t *SelectTable) MarshalState(m maps.Setter) {
	m.Set("selId", t.selectedID)
}

func (t *SelectTable) UnmarshalState(m maps.Loader) {
	if v, ok := m.Load("selId"); ok {
		if id, ok2 := v.(string); ok2 {
			t.selectedID = id
		}
	}
}

func (t *SelectTable) Serialize(e page.Encoder) {
	t.Table.Serialize(e)
	if err := e.Encode(t.selectedID); err != nil {
		panic(err)
	}
	if err := e.Encode(t.reselectable); err != nil {
		panic(err)
	}
}
func (t *SelectTable) Deserialize(dec page.Decoder) {
	t.Table.Deserialize(dec)

	if err := dec.Decode(&t.selectedID); err != nil {
		panic(err)
	}
	if err := dec.Decode(&t.reselectable); err != nil {
		panic(err)
	}
}

// SelectTableCreator is the initialization structure for declarative creation of tables
type SelectTableCreator struct {

	// ID is the control id
	ID               string
	// Caption is the content of the caption tag, and can either be a string, or a data pager
	Caption          interface{}
	// HideIfEmpty will hide the table completely if it has no data. Otherwise, the table and headers will be shown, but no data rows
	HideIfEmpty      bool
	// HeaderRowCount is the number of header rows. You must set this to at least 1 to show header rows.
	HeaderRowCount   int
	// FooterRowCount is the number of footer rows.
	FooterRowCount   int
	// RowStyler returns the attributes to be used in a cell.
	RowStyler        TableRowAttributer
	// RowStylerID is a control id for the control that will be the RowStyler of the table.
	RowStylerID      string
	// HeaderRowStyler returns the attributes to be used in a header cell.
	HeaderRowStyler  TableHeaderRowAttributer
	// HeaderRowStylerID is a control id for the control that will be the HeaderRowStyler of the table.
	HeaderRowStylerID  string
	// FooterRowStyler returns the attributes to be used in a footer cell. It can be either a control id or a TableFooterRowAttributer.
	FooterRowStyler  TableFooterRowAttributer
	// FooterRowStylerID is a control id for the control that will be the FooterRowStyler of the table.
	FooterRowStylerID  string
	// Columns are the column creators that will add columns to the table
	Columns          []ColumnCreator
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// Data is the actual data for the table, and should be a slice of objects
	Data             interface{}
	// Sortable will make the table sortable
	Sortable         bool
	// SortHistoryLimit will set how many columns deep we will remember the sorting for multi-level sorts
	SortHistoryLimit int
	page.ControlOptions
	// OnRowSelected is the action to take when the row is selected
	OnRowSelected    action.ActionI
	// SelectedID is the row id that will start as the selection
	SelectedID string
	// Reselectable determines if you will get a select command when the user taps the item that is already selected.
	Reselectable bool
	// SaveState will cause the table to remember the selection
	SaveState bool
}



// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c SelectTableCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewSelectTable(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator. You do not normally need to call this.
func (c SelectTableCreator) Init(ctx context.Context, ctrl SelectTableI) {
	sub := TableCreator {
		ID:               c.ID,
		Caption:          c.Caption,
		HideIfEmpty:      c.HideIfEmpty,
		HeaderRowCount:   c.HeaderRowCount,
		FooterRowCount:   c.FooterRowCount,
		RowStyler:        c.RowStyler,
		HeaderRowStyler:  c.HeaderRowStyler,
		FooterRowStyler:  c.FooterRowStyler,
		Columns:          c.Columns,
		DataProvider:     c.DataProvider,
		DataProviderID:   c.DataProviderID,
		Data:             c.Data,
		Sortable:         c.Sortable,
		SortHistoryLimit: c.SortHistoryLimit,
		ControlOptions:   c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
	if c.SelectedID != "" {
		ctrl.SetSelectedID(c.SelectedID)
	}
	if c.OnRowSelected != nil {
		ctrl.On(event.RowSelected(), c.OnRowSelected)
	}
	ctrl.SetReselectable(c.Reselectable)
	if c.SaveState { // will override the initial SelectedID setting if true
		ctrl.SaveState(ctx, true)
	}
}

// GetSelectTable is a convenience method to return the button with the given id from the page.
func GetSelectTable(c page.ControlI, id string) *SelectTable {
	return c.Page().GetControl(id).(*SelectTable)
}

func init() {
	page.RegisterControl(&SelectTable{})
}
