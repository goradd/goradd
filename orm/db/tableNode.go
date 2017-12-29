package db

import (
	"log"
	"strings"
)

// A TableNode is a representation of the top level of a chain of nodes that point to a particular field in a query, even after
// aliases and joins are taken into account.

type TableNodeI interface {
	NodeI
	SelectNodes_() []*ColumnNode
	PrimaryKeyNode_() *ColumnNode
	EmbeddedNode_() NodeI
}

type TableNode struct {
	Node

	// Which database in the global list of databases does the node belong to
	dbKey			string
	// Name of table in the database we point to
	dbTable       string
	// The name of the function used to access the property as a node or ORM item
	goName			string
}


// NewTableNode creates a table node, which is always the starting node in a node chain
func NewTableNode(dbKey string, dbName string, goName string) *TableNode {
	return &TableNode {
		dbKey:   dbKey,
		dbTable: dbName,
		goName:  goName,
	}
}

func (n *TableNode) nodeType() NodeType {
	return TABLE_NODE
}

func (n *TableNode) tableName() string {
	return n.dbTable
}

func (n *TableNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == TABLE_NODE {
		cn := n2.(TableNodeI).EmbeddedNode_().(*TableNode)
		return cn.dbTable == n.dbTable && cn.dbKey == n.dbKey
	}

	return false
}

func (n *TableNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Table: " + n.dbTable + " AS " + n.GetAlias())
}


// Return the name as a captialized object name
func (n *TableNode) objectName() string {
	return n.goName
}