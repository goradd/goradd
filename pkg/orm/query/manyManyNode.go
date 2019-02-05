package query

import (
	"log"
	"strings"
)

// A ManyManyNode is an association node that links one table to another table with a many-to-many relationship.
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

// NewManyManyNode  is used internally by the framework to return a new ManyMany node.
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
	return ManyManyNodeType
}

// Expand tells this node to create multiple original objects with a single link for each joined item, rather than to create one original with an array of joined items
func (n *ManyManyNode) Expand() {
	n.isArray = false
}

// isExpanded reports whether this node is creating a new object for each joined item (true), or creating an array of
// joined items (false).
func (n *ManyManyNode) isExpanded() bool {
	return !n.isArray
}

// Equals is used internally by the framework to test if the node is the same as another node.
func (n *ManyManyNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == ManyManyNodeType {
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

// ManyManyNodeIsArray is used internally by the framework to return whether the node creates an array, or just a link to a single item.
func ManyManyNodeIsArray(n *ManyManyNode) bool {
	return n.isArray
}

// ManyManyNodeIsTypeTable is used internally by the framework to return whether the node points to a type table
func ManyManyNodeIsTypeTable(n *ManyManyNode) bool {
	return n.isTypeTable
}

// ManyManyNodeRefTable is used internally by the framework to return the table name on the other end of the link
func ManyManyNodeRefTable(n *ManyManyNode) string {
	return n.refTable
}

// ManyManyNodeRefColumn is used internally by the framework to return the column name on the other end of the link
func ManyManyNodeRefColumn(n *ManyManyNode) string {
	return n.refColumn
}

// ManyManyNodeDbTable is used internally by the framework to return the table name of the table the node belongs to
func ManyManyNodeDbTable(n *ManyManyNode) string {
	return n.dbTable
}

// ManyManyNodeDbColumn is used internally by the framework to return the column name in the table the node belongs to
func ManyManyNodeDbColumn(n *ManyManyNode) string {
	return n.dbColumn
}
