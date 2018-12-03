package column

import (
	"github.com/spekary/goradd/pkg/page/control"
	"goradd-project/override/control_base"
)


type ProxyCellTexter interface {
	control.ProxyI
	CellTexter
}

// ProxyColumn is a table column that prints the output of a Proxy control.
// To use it, you must define your own proxy control that also has a CellText method
// attached to it, so that it satisfies the ProxyCellTexter interface above
type ProxyColumn struct {
	control.ColumnBase
	Proxy ProxyCellTexter
}

// NewProxyColumn creates a new column with a custom cell texter.
//
// Set isHtml to true to indicate that the cell texter is returning html and not plain text.
func NewProxyColumn(proxy ProxyCellTexter) *ProxyColumn {
	i := ProxyColumn{}
	i.Init(proxy)
	return &i
}

func (c *ProxyColumn) Init(proxy ProxyCellTexter) {
	c.ColumnBase.Init(c)
	c.ColumnBase.SetIsHtml(true)
	c.SetCellTexter(proxy)
}
