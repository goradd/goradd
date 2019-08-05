package column

import (
	"github.com/goradd/goradd/pkg/page/control"
)

// CustomColumn is a table column that you can customize any way you want. You simply give it a CellTexter, and return
// the text from the cell texter. One convenient way to use this is to define a CellText function on the
// parent object and pass it as the CellTexter. If your CellTexter is going to output html instead of raw text, call
// SetIsHtml() on the column after creating it.
type CustomColumn struct {
	control.ColumnBase
}

// NewCustomColumn creates a new column with a custom cell texter.
func NewCustomColumn(texter CellTexter) *CustomColumn {
	i := CustomColumn{}
	i.Init(texter)
	return &i
}

func (c *CustomColumn) Init(texter CellTexter) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(texter)
}

// Just make ColumnOptions completely available
type CustomColumnCreator control.ColumnOptions

func (c CustomColumnCreator) Create(parent control.TableI) control.ColumnI {
	col := NewCustomColumn(GetCellTexter(parent, c.CellTexterID))
	col.ApplyOptions(parent, control.ColumnOptions(c))
	return col
}