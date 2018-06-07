package query

import (
	"log"
	"strings"
)

// ManyManyNode creates an association node.
// Some of the columns have overloaded meanings depending on SQL or NoSQL mode.
type ManyManyNode struct {
	Node

	// Which database in the global list of databases does the node belong to
	dbKey string
	// NoSQL: The originating table. SQL: The association table
	dbTable string
	// NoSQL: The table storing the array of ids on the other end. SQL: the table in the association table pointing towards us.
	dbColumn string
	// Property in the original object used to ref to this object or node.
	goPropName string

	// NoSQL & SQL: The table we are joining to
	refTable string
	// NoSQL: table point backwards to us. SQL: Column in association table pointing forwards to refTable
	refColumn string
	// Are we expanding as an array, or one item at a time.
	isArray bool
	// Is this pointing to a type table item?
	isTypeTable bool
}

func NewManyManyNode(
	dbKey string,
	// NoSQL: The originating table. SQL: The association table
	dbTable string,
	// NoSQL: The table storing the array of ids on the other end. SQL: the table in the association table pointing towards us.
	dbColumn string,
	// Property in the original object used to ref to this object or node.
	goName string,
	// NoSQL & SQL: The table we are joining to
	refTableName string,
	// NoSQL: table point backwards to us. SQL: Column in association table pointing forwards to refTable
	refColumn string,
	// Are we pointing to a type table
	isType bool,
) *ManyManyNode {
	n := &ManyManyNode{
		dbKey:       dbKey,
		dbTable:     dbTable,
		dbColumn:    dbColumn,
		goPropName:  goName,
		refTable:    refTableName,
		refColumn:   refColumn,
		isArray:     true,
		isTypeTable: isType,
	}
	return n
}

func (n *ManyManyNode) nodeType() NodeType {
	return MANYMANY_NODE
}

func (n *ManyManyNode) Expand() {
	n.isArray = false
}

func (n *ManyManyNode) isExpanded() bool {
	return !n.isArray
}

func (n *ManyManyNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == MANYMANY_NODE {
		cn := n2.(TableNodeI).EmbeddedNode_().(*ManyManyNode)
		return cn.dbTable == n.dbTable &&
			cn.goPropName == n.goPropName &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)

	}
	return false
}

func (n *ManyManyNode) setCondition(condition NodeI) {
	n.condition = condition
}

func (n *ManyManyNode) getCondition() NodeI {
	return n.condition
}

func (n *ManyManyNode) tableName() string {
	return n.refTable
}

func (n *ManyManyNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "MM: " + n.dbTable + "." + n.dbColumn + "." + n.refTable + "." + n.refColumn + " AS " + n.GetAlias())
}

// Return the name as a captialized object name
func (n *ManyManyNode) goName() string {
	return n.goPropName
}

func ManyManyNodeIsArray(n *ManyManyNode) bool {
	return n.isArray
}

func ManyManyNodeIsTypeTable(n *ManyManyNode) bool {
	return n.isTypeTable
}

func ManyManyNodeRefTable(n *ManyManyNode) string {
	return n.refTable
}

func ManyManyNodeRefColumn(n *ManyManyNode) string {
	return n.refColumn
}

func ManyManyNodeDbTable(n *ManyManyNode) string {
	return n.dbTable
}

func ManyManyNodeDbColumn(n *ManyManyNode) string {
	return n.dbColumn
}
