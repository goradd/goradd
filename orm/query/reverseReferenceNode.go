package query

import (
	"log"
	"strings"
)

// ReverseReferenceNode creates a reverse reference node representing a one to many relationship. This is the many side of that relationship.
type ReverseReferenceNode struct {
	Node

	// Which database in the global list of databases does the node belong to
	dbKey string
	// table which has the reverse reference
	dbTable string
	// NoSQL only. The table containing an array of items we are pointing to.
	dbColumn string
	// Property we are using to refer to the many side of the relationship
	goPropName string

	// Is this pointing to a type table item?
	isTypeTable bool

	// The table containing the pointer back to us
	refTable string
	// The table that is the foreign key pointing back to us.
	refColumn string

	// True to create new objects for each joined item, or false to create an array of joined objects here.
	isArray bool
}

func NewReverseReferenceNode(
	dbKey string,
	// table which has the reverse reference
	dbTable string,
	// NoSQL: the table containing an array of items we are pointing to. SQL: The primary key of this table.
	dbColumn string,
	// Property we are using to refer to the many side of the relationship
	goName string,
	// The table containing the pointer back to us
	refTable string,
	// The table that is the foreign key pointing back to us.
	refColumn string,
	isArray bool,
) *ReverseReferenceNode {
	n := &ReverseReferenceNode{
		dbKey:      dbKey,
		dbTable:    dbTable,
		dbColumn:   dbColumn,
		goPropName: goName,
		refTable:   refTable,
		refColumn:  refColumn,
		isArray:    isArray,
	}
	return n
}

func (n *ReverseReferenceNode) nodeType() NodeType {
	return REVERSE_REFERENCE_NODE
}

func (n *ReverseReferenceNode) Expand() {
	n.isArray = false
}

func (n *ReverseReferenceNode) isExpanded() bool {
	return !n.isArray
}

func (n *ReverseReferenceNode) IsArray() bool {
	return n.isArray
}

func (n *ReverseReferenceNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == REVERSE_REFERENCE_NODE {
		cn := n2.(TableNodeI).EmbeddedNode_().(*ReverseReferenceNode)
		return cn.dbTable == n.dbTable &&
			cn.refTable == n.refTable &&
			cn.refColumn == n.refColumn &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)
	}

	return false
}

func (n *ReverseReferenceNode) tableName() string {
	return n.refTable
}

func (n *ReverseReferenceNode) setCondition(condition NodeI) {
	n.condition = condition
}

func (n *ReverseReferenceNode) getCondition() NodeI {
	return n.condition
}

func (n *ReverseReferenceNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "RR: " + n.dbTable + "." + n.refTable + "." + n.refColumn + " AS " + n.GetAlias())
}

// Return the name as a captialized object name
func (n *ReverseReferenceNode) goName() string {
	return n.goPropName
}

func ReverseReferenceNodeIsArray(n *ReverseReferenceNode) bool {
	return n.isArray
}

func ReverseReferenceNodeRefTable(n *ReverseReferenceNode) string {
	return n.refTable
}

func ReverseReferenceNodeRefColumn(n *ReverseReferenceNode) string {
	return n.refColumn
}

func ReverseReferenceNodeDbColumnName(n *ReverseReferenceNode) string {
	return n.dbColumn
}
