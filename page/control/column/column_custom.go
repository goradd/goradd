package column

import (
	"goradd/override/control_base"
)

// CustomColumn is a table column that you can customize any way you want. You simply give it a CellTexter, and return
// the text or html from the cell texter.
type CustomColumn struct {
	control_base.ColumnBase
}

// NewCustomColumn creates a new column with a custom cell texter.
//
// Set isHtml to true to indicate that the cell texter is returning html and not plain text.
func NewCustomColumn(texter CellTexter, isHtml bool) *CustomColumn {
	i := CustomColumn{}
	i.Init(texter, isHtml)
	return &i
}

func (c *CustomColumn) Init(texter CellTexter, isHtml bool) {
	c.ColumnBase.Init(c)
	c.ColumnBase.SetIsHtml(isHtml)
	c.SetCellTexter(texter)
}
