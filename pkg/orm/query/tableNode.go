package query

import (
	"log"
	"strings"
)

// TableNodeI is the interface that all table nodes must implement. TableNodes are create by the code generation
// process, one for each table in the database.
type TableNodeI interface {
	ReferenceNodeI
	SelectNodes_() []*ColumnNode
	PrimaryKeyNode_() *ColumnNode
	EmbeddedNode_() NodeI
	Copy_() NodeI
}

// A TableNode is a representation of the top level of a chain of nodes that point to a particular field in a query, even after
// aliases and joins are taken into account. TableNodes are create by the code generation
// process, one for each table in the database.
type TableNode struct {
	nodeAlias
	nodeLink
	// Which database in the global list of databases does the node belong to
	dbKey string
	// Name of table in the database we point to
	dbTable string
	// The name of the function used to access the property as a node or ORM item
	goPropName string
}

// NewTableNode creates a table node, which is always the starting node in a node chain
func NewTableNode(dbKey string, dbName string, goName string) *TableNode {
	return &TableNode{
		dbKey:      dbKey,
		dbTable:    dbName,
		goPropName: goName,
	}
}

func (n *TableNode) copy() NodeI {
	return &TableNode{
		dbKey:      n.dbKey,
		dbTable:    n.dbTable,
		goPropName: n.goPropName,
		nodeAlias:  nodeAlias{n.alias},
	}
}

func (n *TableNode) tableName() string {
	return n.dbTable
}

func (n *TableNode) goName() string {
	return n.goPropName
}

func (n *TableNode) Equals(n2 NodeI) bool {
	if tn, ok := n2.(TableNodeI); !ok {
		return false
	} else if cn, ok := tn.EmbeddedNode_().(*TableNode); !ok {
		return false
	} else {
		return cn.dbTable == n.dbTable &&
			cn.dbKey == n.dbKey &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)
	}
}

func (n *TableNode) Expand() {
	panic("you cannot expand a TableNode")
}

func (n *TableNode) isExpanded() bool {
	return false
}

func (n *TableNode) getCondition() NodeI {
	return nil
}

func (n *TableNode) setCondition(c NodeI) {
	panic("you cannot set a condition on a TableNode")
}

func (n *TableNode) nodeType() NodeType {
	return TableNodeType
}

func (n *TableNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Table: " + n.dbTable)
}
