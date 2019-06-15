package column

import (
	"context"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page/control"
	"reflect"
)

// NodeColumn is a column that uses a query.NodeI to get its text out of data that is coming from the ORM.
// Create it with NewNodeColumn
type NodeColumn struct {
	control.ColumnBase
}

// NewNodeColumn creates a table column that uses a query.NodeI object to get its text out of an ORM object.
// node should point to data that is preloaded in the ORM object. If the node points to a Date, Time or DateTime type
// of column, you MUST specify a time format by calling SetTimeFormat.
func NewNodeColumn(node query.NodeI) *NodeColumn {
	i := NodeColumn{}
	i.Init(node)
	return &i
}

func (c *NodeColumn) Init(node query.NodeI) {
	if node == nil {
		panic("node is required")
	}
	c.ColumnBase.Init(c)
	c.SetCellTexter(&NodeTexter{Node: node})
	c.SetTitle(query.NodeGoName(node))
}

func (c *NodeColumn) GetNode() query.NodeI {
	return c.CellTexter().(*NodeTexter).Node
}

// SetFormat sets the format string of the node column.
func (c *NodeColumn) SetFormat(format string) *NodeColumn {
	c.CellTexter().(*NodeTexter).Format = format
	return c
}

// SetTimeFormat sets the time format of the string, specifically for a DateTime column.
func (c *NodeColumn) SetTimeFormat(format string) *NodeColumn {
	c.CellTexter().(*NodeTexter).TimeFormat = format
	return c
}

// NodeTexter is used by the NodeColumn to get text out of a database column.
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

func (t NodeTexter) CellText(ctx context.Context, col control.ColumnI, rowNum int, colNum int, data interface{}) string {
	if v, ok := data.(Getter); !ok {
		return ""
	} else {
		n := t.Node
		var names []string

		// walk up the chain of nodes to figure out how to walk down the chain of data
		for {
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
			v2, ok = obj.(Getter)
			if !ok {
				panic("node chain does not match a chain of Getters (forward, reverse and manyMany references")
			}
		}
		s := v2.Get(names[0]) // This should be a column node
		return ApplyFormat(s, t.Format, t.TimeFormat)
	}
}

type NodeGetter interface {
	GetNode() query.NodeI
}

// MakeNodeSlice is a convenience method to convert a slice of columns into a slice of nodes derived from
// those columns. The column slice would typically come from the table's SortColumns method, and the returned
// slice would be passed to the database's OrderBy clause when building a query. Since this is a common use, it
// will also add sort info to the nodes.
func MakeNodeSlice(columns []control.ColumnI) []query.NodeI {
	var nodes []query.NodeI
	for _, c := range columns {
		if getter, ok := c.(NodeGetter); ok {
			node := getter.GetNode()
			if nodeSorter, ok := node.(query.NodeSorter); ok {
				switch c.SortDirection() {
				case control.SortAscending:
					nodeSorter.Ascending()
				case control.SortDescending:
					nodeSorter.Descending()
				}
			}
			nodes = append(nodes, node)
		} else {
			panic("Column is not a NodeGetter.")
		}
	}
	return nodes
}
