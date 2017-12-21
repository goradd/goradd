package db

import (
	"log"
	"strings"
)

// ManyManyNode creates an association node.
// Some of the columns have overloaded meanings depending on SQL or NoSQL mode.
type ManyManyNode struct {
	Node

	// Which database in the global list of databases does the node belong to
	dbKey			string
	// NoSQL: The originating table. SQL: The association table
	dbTable       string
	// NoSQL: The column storing the array of ids on the other end. SQL: the column in the association table pointing towards us.
	dbColumn       string
	// Property in the original object used to ref to this object or node.
	goName			string

	// NoSQL & SQL: The table we are joining to
	refTable string
	// NoSQL: column point backwards to us. SQL: Column in association table pointing forwards to refTable
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
// NoSQL: The column storing the array of ids on the other end. SQL: the column in the association table pointing towards us.
	dbColumn string,
// Property in the original object used to ref to this object or node.
	goName string,
// NoSQL & SQL: The table we are joining to
	refTableName string,
// NoSQL: column point backwards to us. SQL: Column in association table pointing forwards to refTable
	refColumn string,
// Are we pointing to a type table
	isType bool,
) *ManyManyNode {
	n:= &ManyManyNode {
		dbKey:       dbKey,
		dbTable:     dbTable,
		dbColumn:    dbColumn,
		goName:      goName,
		refTable:    refTableName,
		refColumn:   refColumn,
		isArray: 	 true,
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
		return cn.dbTable == n.dbTable && cn.goName == n.goName
	}
	return false
}

func (n *ManyManyNode) setConditions(conditions []NodeI){
	n.conditions = conditions
}

func (n *ManyManyNode) getConditions() []NodeI {
	return n.conditions
}

func (n *ManyManyNode) tableName() string {
	return n.refTable
}

func (n *ManyManyNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "MM: " + n.dbTable + "." + n.dbColumn + "." + n.refTable + "." + n.refColumn)
}


// Return the name as a captialized object name
func (n *ManyManyNode) objectName() string {
	return n.goName
}