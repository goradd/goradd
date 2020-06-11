package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control"
)

// TexterColumn is a table column that lets you use a CellTexter to specify the content of the cells in the column.
// One convenient way to use this is to define a CellText function on the
// parent object and pass it as the CellTexter. If your CellTexter is going to output html instead of raw text, call
// SetIsHtml() on the column after creating it.
type TexterColumn struct {
	control.ColumnBase
}

// NewTexterColumn creates a new column with a custom cell texter.
func NewTexterColumn(texter CellTexter) *TexterColumn {
	i := TexterColumn{}
	i.Init(texter)
	return &i
}

func (c *TexterColumn) Init(texter CellTexter) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(texter)
}

// TexterColumnCreator creates a column that uses a CellTexter to get the content of each cell.
type TexterColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Texter is a CellTexter that returns the text that should go in each cell.
	Texter CellTexter
	// TexterID is the control id of a CellTexter.
	TexterID string
	// Title is the title at the top of the column
	Title string
	// Sortable makes the column display sort arrows in the header
	Sortable bool
	// IsHtml indicates that the texter is producing HTML rather than text that should be escaped.
	IsHtml bool
	control.ColumnOptions
}

func (c TexterColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	if c.Texter == nil && c.TexterID == "" {
		panic("a Texter is required")
	}
	var texter CellTexter
	if c.Texter != nil {
		texter = c.Texter
	} else {
		texter = GetCellTexter(parent, c.TexterID)
	}

	col := NewTexterColumn(texter)
	if c.ID != "" {
		col.SetID(c.ID)
	}
	col.SetTitle(c.Title)
	if c.Sortable {
		col.SetSortable()
	}
	if c.IsHtml {
		col.SetIsHtml(true)
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}

func init() {
	control.RegisterColumn(TexterColumn{})
}