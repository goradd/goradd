package column

import (
	"context"
	"github.com/goradd/goradd/pkg/page/control/table"

	"github.com/goradd/goradd/pkg/iface"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/page"
)

// NodeColumn is a column that uses a query.NodeI to get its text out of data that is coming from the ORM.
// Create it with NewNodeColumn
type NodeColumn struct {
	table.ColumnBase
	node query.NodeI
}

// NewNodeColumn creates a table column that uses a query.NodeI object to get its text out of an ORM object.
// node should point to data that is preloaded in the ORM object. If the node points to a Date, Time or DateTime type
// of column, you MUST specify a time format by calling SetTimeFormat.
func NewNodeColumn(node query.NodeI) *NodeColumn {
	i := NodeColumn{}
	i.Init(node)
	i.node = node
	return &i
}

func (c *NodeColumn) Init(node query.NodeI) {
	if node == nil {
		panic("node is required")
	}
	c.ColumnBase.Init(c)
	c.SetTitle(query.NodeGoName(node))
}

func (c *NodeColumn) GetNode() query.NodeI {
	return c.node
}

func (c *NodeColumn) CellData(_ context.Context, _ int, _ int, data interface{}) interface{} {
	if v, ok := data.(Getter); !ok {
		return ""
	} else {
		n := c.node
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
			if iface.IsNil(obj) {
				return nil // Either an object was not joined, or the joined object does not exist in the database.
			}
			v2, ok = obj.(Getter)
			if !ok {
				panic("node chain does not match a chain of Getters")
			}
		}
		s := v2.Get(names[0]) // This should be a column node
		return s
	}
}

func (c *NodeColumn) Serialize(e page.Encoder) {
	c.ColumnBase.Serialize(e)
	if err := e.Encode(&c.node); err != nil {
		panic(err)
	}
	return
}

func (c *NodeColumn) Deserialize(dec page.Decoder) {
	c.ColumnBase.Deserialize(dec)
	if err := dec.Decode(&c.node); err != nil {
		panic(err)
	}
}

type NodeGetter interface {
	GetNode() query.NodeI
}

// MakeNodeSlice is a convenience method to convert a slice of columns into a slice of nodes derived from
// those columns. The column slice would typically come from the table's SortColumns method, and the returned
// slice would be passed to the database's OrderBy clause when building a query. Since this is a common use, it
// will also add sort info to the nodes.
func MakeNodeSlice(columns []table.ColumnI) []query.NodeI {
	var nodes []query.NodeI
	for _, c := range columns {
		if getter, ok := c.(NodeGetter); ok {
			node := getter.GetNode()
			if node != nil {
				if nodeSorter, ok2 := node.(query.NodeSorter); ok2 {
					switch c.SortDirection() {
					case table.SortAscending:
						nodeSorter.Ascending()
					case table.SortDescending:
						nodeSorter.Descending()
					}
				}
				nodes = append(nodes, node)
			} else {
				panic("Column does not have a sort node.")
			}
		} else {
			panic("Column is not a NodeGetter.")
		}
	}
	return nodes
}

// NodeColumnCreator creates a column that treats each row as data from the ORM, and gets to that data using
// a database Node.
type NodeColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Node is a database node generated by the code generator
	Node query.NodeI
	// Title is the title of the column that will appear in the header
	Title string
	// Sortable makes the column display sort arrows in the header
	// Deprecated: Use SortDirection instead
	Sortable bool
	// SortDirection sets the initial sorting direction of the column, and will make the column sortable
	// By default, the column is not sortable.
	SortDirection table.SortDirection
	// IsHtml indicates that the texter is producing HTML rather than text that should be escaped.
	table.ColumnOptions
}

func (c NodeColumnCreator) Create(ctx context.Context, parent table.TableI) table.ColumnI {
	col := NewNodeColumn(c.Node)
	if c.ID != "" {
		col.SetID(c.ID)
	}
	col.SetTitle(c.Title)
	if c.Sortable {
		col.SetSortable()
	}
	if c.SortDirection != table.NotSortable {
		col.SetSortDirection(c.SortDirection)
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}

func init() {
	table.RegisterColumn(NodeColumn{})
}
