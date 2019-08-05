package control

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
)

type PaginatedTableI interface {
	TableEmbedder
	PaginatedControlI
}

type PaginatedTable struct {
	Table
	PaginatedControl
}

func NewPaginatedTable(parent page.ControlI, id string) *PaginatedTable {
	t := &PaginatedTable{}
	t.Init(t, parent, id)
	return t
}

func (t *PaginatedTable) Init(self page.ControlI, parent page.ControlI, id string) {
	t.Table.Init(self, parent, id)
	t.PaginatedControl.SetPageSize(0) // use the application default
}

// PaginatedTableCreator is the initialization structure for declarative creation of tables
type PaginatedTableCreator struct {
	// ID is the control id
	ID string
	RenderColumnTags bool
	Caption interface{} // string or paginator
	HideIfEmpty bool
	HeaderRowCount int
	FooterRowCount int
	RowStyler string
	HeaderRowStyler string
	FooterRowStyler string
	Columns []ColumnCreator
	PageSize int
	DataProvider string
	SaveState bool
	page.ControlOptions
}



// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c PaginatedTableCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewPaginatedTable(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator. You do not normally need to call this.
func (c PaginatedTableCreator) Init(ctx context.Context, ctrl PaginatedTableI) {
	sub := TableCreator {
		ID: c.ID,
		RenderColumnTags: c.RenderColumnTags,
		Caption: c.Caption,
		HideIfEmpty: c.HideIfEmpty,
		HeaderRowCount: c.HeaderRowCount,
		FooterRowCount: c.FooterRowCount,
		RowStyler: c.RowStyler,
		HeaderRowStyler: c.HeaderRowStyler,
		FooterRowStyler: c.FooterRowStyler,
		Columns: c.Columns,
		DataProvider: c.DataProvider,
		ControlOptions: c.ControlOptions,
	}
	sub.Init(ctx, ctrl)
	if c.PageSize != 0 {
		ctrl.SetPageSize(c.PageSize)
	}
	if c.SaveState {
		ctrl.SaveState(ctx, true)
	}
}

// GetPaginatedTable is a convenience method to return the button with the given id from the page.
func GetPaginatedTable(c page.ControlI, id string) *PaginatedTable {
	return c.Page().GetControl(id).(*PaginatedTable)
}
