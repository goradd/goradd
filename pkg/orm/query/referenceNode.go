package query

import (
	"log"
	"strings"
)

type ReferenceNodeI interface {
	NodeI
	Aliaser
	conditioner
	nodeLinkI
	Expander
}

// A ReferenceNode is a forward-pointing foreign key relationship, and can define a one-to-one or
// one-to-many relationship, depending on whether it is unique. If the other side of the relationship is
// not a type table, then the other table will have a matching ReverseReferenceNode.
type ReferenceNode struct {
	nodeAlias
	nodeCondition
	nodeLink
	// Which database in the global list of databases does the node belong to
	dbKey string
	// Name of table in the database we point to
	dbTable string
	// The name of the table that is the foreign key
	dbColumn string
	// The name of the table related to this reference
	goColumnName string
	// The name of the function used to access the property as a node or ORM item
	goPropName string
	// The name of the variable in the model structure used to hold the object
	goVarName string

	// Is this pointing to a type table item?
	isTypeTable bool

	// The name of the table we are joining to
	refTable string
	// If a forward reference and NoSQL, the name of the table that will contain the reference or references backwards to us. If SQL, the Pk of the RefTable
	refColumn string
}

// NewReferenceNode creates a forward reference node.
func NewReferenceNode(
	dbKey string,
	dbTableName string,
	dbColumnName string,
	goColumnName string,
	goName string,
	refTableName string,
	refColumn string, // only used in NoSQL situation
	isType bool,
) *ReferenceNode {
	n := &ReferenceNode{
		dbKey:        dbKey,
		dbTable:      dbTableName,
		dbColumn:     dbColumnName,
		goColumnName: goColumnName,
		goPropName:   goName,
		refTable:     refTableName,
		refColumn:    refColumn,
		isTypeTable:  isType,
	}
	return n
}

func (n *ReferenceNode) copy() NodeI {
	ret := &ReferenceNode{
		dbKey:         n.dbKey,
		dbTable:       n.dbTable,
		dbColumn:      n.dbColumn,
		goColumnName:  n.goColumnName,
		goPropName:    n.goPropName,
		refTable:      n.refTable,
		refColumn:     n.refColumn,
		isTypeTable:   n.isTypeTable,
		nodeAlias:     nodeAlias{n.alias},
		nodeCondition: nodeCondition{n.condition},
	}
	return ret
}

// Equals is used internally by the framework to determine if two nodes are equal.
func (n *ReferenceNode) Equals(n2 NodeI) bool {
	if tn, ok := n2.(TableNodeI); !ok {
		return false
	} else if cn, ok := tn.EmbeddedNode_().(*ReferenceNode); !ok {
		return false
	} else {
		return cn.dbTable == n.dbTable &&
			cn.goPropName == n.goPropName &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)
	}
}

func (n *ReferenceNode) nodeType() NodeType {
	return ReferenceNodeType
}

func (n *ReferenceNode) tableName() string {
	return n.refTable
}

func (n *ReferenceNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "R: " + n.dbTable + "." + n.dbColumn + "." + n.refTable + " AS " + n.GetAlias())
}

// Return the name as a capitalized object name
func (n *ReferenceNode) goName() string {
	return n.goPropName
}

// Return a column node for the foreign key that represents the reference to the other table
func (n *ReferenceNode) relatedColumnNode() *ColumnNode {
	var colType GoColumnType
	if n.isTypeTable {
		colType = ColTypeUnsigned
	} else {
		colType = ColTypeString
	}

	n2 := NewColumnNode(n.dbKey, n.dbTable, n.dbColumn, n.goColumnName, colType, false)
	SetParentNode(n2, n.getParent())
	return n2
}

func (n *ReferenceNode) Expand() {
	panic("you cannot expand on a reference node, only reverse reference and many-many reference")
}

func (n *ReferenceNode) isExpanded() bool {
	return false
}

// RelatedColumnNode is used internally by the framework to create a new node for the other side of the relationship.
func RelatedColumnNode(n NodeI) NodeI {
	if tn, _ := n.(TableNodeI); tn != nil {
		if rn, _ := tn.EmbeddedNode_().(*ReferenceNode); rn != nil {
			return rn.relatedColumnNode()
		}
	}
	return nil
}

// ReferenceNodeRefTable is used internally by the framework to get the table name for the other side of the relationship.
func ReferenceNodeRefTable(n *ReferenceNode) string {
	return n.refTable
}

// ReferenceNodeRefColumn is used internally by the framework to get the column name for the other side of the relationship.
func ReferenceNodeRefColumn(n *ReferenceNode) string {
	return n.refColumn
}

// ReferenceNodeDbColumnName is used internally by the framework to get the column name for this side of the relationship.
func ReferenceNodeDbColumnName(n *ReferenceNode) string {
	return n.dbColumn
}
