package column

import (
	"context"
	"goradd-project/override/control_base"
	"github.com/spekary/goradd/page/control/control_base/table"
	"github.com/spekary/goradd/orm/query"
	"reflect"
)

// NodeColumn is a column that uses the Getter interface to get the text out of columns. The data therefore should be
// a slice of objects that implement the Getter interface.
type NodeColumn struct {
	control_base.ColumnBase
	node query.NodeI
}

func NewNodeColumn(node query.NodeI, format ...string) *NodeColumn {
	i := NodeColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(node, f, "")
	return &i
}

func NewDateNodeColumn(node query.NodeI, timeFormat string, format ...string) *NodeColumn {
	i := NodeColumn{}
	var f string
	if len(format) > 0 {
		f = format[0]
	}
	i.Init(node, f, timeFormat)
	return &i
}

func (c *NodeColumn) Init(node query.NodeI, format string, timeFormat string) {
	c.ColumnBase.Init(c)
	c.SetCellTexter(&NodeTexter{Node: node, Format: format, TimeFormat: timeFormat})
	c.SetTitle(query.NodeGoName(node))
}

func (c *NodeColumn) SetFormat(format string) *NodeColumn {
	c.CellTexter().(*NodeTexter).Format = format
	return c
}

func (c *NodeColumn) SetTimeFormat(format string) *NodeColumn {
	c.CellTexter().(*NodeTexter).TimeFormat = format
	return c
}

// GetterTexter lets you get items out of map like objects using the Getter interface.
type NodeTexter struct {
	// Key is the key to use when calling the Get function on the object.
	Node query.NodeI
	// Format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	Format string
	// TimeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	TimeFormat string
}

func (t NodeTexter) CellText(ctx context.Context, col table.ColumnI, rowNum int, colNum int, data interface{}) string {
	if v,ok := data.(Getter); !ok {
		return ""
	} else {
		n := t.Node
		var names []string

		// walk up the chain of nodes to figure out how to walk down the chain of data
		for  {
			name := query.NodeGoName(n)
			if name == "" {
				break
			}
			names = append(names, name)
			n = query.ParentNode(n)
			if n == nil {
				break
			}
		}
		if len(names) < 2 {
			panic("bad node passed to the column_node column. These nodes must start with a table node and end with a column node")
		}
		var i int
		v2 := v
		for i = len(names) - 2; i > 0; i-- {
			obj := v2.Get(names[i])
			if obj == nil || reflect.ValueOf(obj).IsNil() {
				panic("database object has not loaded the items referred to in the node. Make sure you are joining the correct tables")
			}
			v2,ok = obj.(Getter)
			if !ok {
				panic ("node chain does not match a chain of Getters (forward, reverse and manyMany references")
			}
		}
		s := v2.Get(names[0]) // This should be a column node
		return ApplyFormat(s, t.Format, t.TimeFormat)
	}
}
