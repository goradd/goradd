package control

import (
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/page/event"
)

// PrimaryKeyer is an interface that is often implemented by model objects.
type PrimaryKeyer interface {
	PrimaryKey() string
}

type SelectTableI interface {
	TableI
	SetSelectedID(id string) SelectTableI
}

// SelectTable is a table that is row selectable. To detect a row selection, trigger on event.RowSelected
type SelectTable struct {
	Table
	selectedID string
}

func NewSelectTable(parent page.ControlI, id string) *SelectTable {
	t := &SelectTable{}
	t.Init(t, parent, id)
	return t
}

func (t *SelectTable) Init(self page.ControlI, parent page.ControlI, id string) {
	t.Table.Init(self, parent, id)
	t.ParentForm().AddJQueryUI()
	t.ParentForm().AddJavaScriptFile(config.GoraddAssets() + "/js/goradd-scrollIntoView.js", false, nil)
	t.ParentForm().AddJavaScriptFile(config.GoraddAssets() + "/js/table-select.js", false, nil)
	t.SetAttribute("tabindex", 0); // Make the entire table focusable and selectable. This can be overridden later if needed.
	t.AddClass("gr-clickable-rows")
}

func (t *SelectTable) this() SelectTableI {
	return t.Self.(SelectTableI)
}

func (t *SelectTable) GetRowAttributes(row int, data interface{}) (a *html.Attributes) {
	var id string

	if t.RowStyler() != nil {
		a = t.RowStyler().TableRowAttributes(row, data)
		id = a.Get("id") // styler might be giving us an id
	} else {
		a = html.NewAttributes()
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
		a.SetDataAttribute("id", id)
		// We need an actual id for aria features
		a.SetID(t.ID() + "_" + id)
	} else {
		a.AddClass("nosel")
	}

	a.Set("role", "option")
	return a
}

func (t *SelectTable) ΩDrawingAttributes() *html.Attributes {
	a := t.Table.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "selecttable")
	a.Set("role", "listbox")
	a.SetDataAttribute("grWidget", "goradd.selectTable")
	if t.selectedID != "" {
		a.SetDataAttribute("grOptSelectedId", t.selectedID)
	}
	return a
}

func (t *SelectTable) ΩUpdateFormValues(ctx *page.Context) {
	if data := ctx.CustomControlValue(t.ID(), "selectedId"); data != nil {
		t.selectedID = fmt.Sprintf("%v", data)
	}
}

func (t *SelectTable) SelectedID() string {
	return t.selectedID
}

func (t *SelectTable) SetSelectedID(id string) SelectTableI {
	t.selectedID = id
	t.ExecuteWidgetFunction("option", "selectedId", id)
	return t.this()
}

func (t *SelectTable) ΩMarshalState(m maps.Setter) {
	m.Set("selId", t.selectedID)
}

func (t *SelectTable) ΩUnmarshalState(m maps.Loader) {
	if v, ok := m.Load("selId"); ok {
		if id, ok := v.(string); ok {
			t.selectedID = id
		}
	}
}


// SelectTableCreator is the initialization structure for declarative creation of tables
type SelectTableCreator struct {

	// ID is the control id
	ID               string
	// HasColTags will make the table render <col> tags
	HasColTags       bool
	// Caption is the content of the caption tag, and can either be a string, or a data pager
	Caption          interface{}
	// HideIfEmpty will hide the table completely if it has no data. Otherwise, the table and headers will be shown, but no data rows
	HideIfEmpty      bool
	// HeaderRowCount is the number of header rows. You must set this to at least 1 to show header rows.
	HeaderRowCount   int
	// FooterRowCount is the number of footer rows.
	FooterRowCount   int
	// RowStyler returns the attributes to be used in a cell. It can be either a control id or a TableRowAttributer.
	RowStyler        interface{}
	// HeaderRowStyler returns the attributes to be used in a header cell. It can be either a control id or a TableHeaderRowAttributer.
	HeaderRowStyler  interface{}
	// FooterRowStyler returns the attributes to be used in a footer cell. It can be either a control id or a TableFooterRowAttributer.
	FooterRowStyler  interface{}
	// Columns are the column creators that will add columns to the table
	Columns          []ColumnCreator
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider data.DataBinder
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
		HasColTags:       c.HasColTags,
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
	if c.SaveState { // will override the initial SelectedID setting if true
		ctrl.SaveState(ctx, true)
	}
	if c.OnRowSelected != nil {
		ctrl.On(event.RowSelected(), c.OnRowSelected)
	}
}

// GetSelectTable is a convenience method to return the button with the given id from the page.
func GetSelectTable(c page.ControlI, id string) *SelectTable {
	return c.Page().GetControl(id).(*SelectTable)
}
