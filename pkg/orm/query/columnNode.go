package query

import (
	"log"
	"strings"
)

// A Column represents a table or field in a database structure, and is the leaf of a node tree or chain.
type ColumnNode struct {
	Node

	// Which database in the global list of databases does the node belong to
	dbKey string
	// Name of table in the database we point to
	dbTable string
	// The name of the table in the database
	dbColumn string
	// The name of the function used to access the property as a node or ORM item
	gName string
	// The go type for the column
	goType GoColumnType
	// Used by OrderBy clauses
	sortDescending bool
	// When in an update or insert operation, stores the value we are attempting to update or insert
	//value interface{}
}

func NewColumnNode(dbKey string, dbTable string, dbName string, goName string, goType GoColumnType) *ColumnNode {
	n := &ColumnNode{
		dbKey:    dbKey,
		dbTable:  dbTable,
		dbColumn: dbName,
		gName:   goName,
		goType:   goType,
	}
	return n
}

// Ascending is used in an OrderBy query builder function to sort the column in ascending order.
func (n *ColumnNode) Ascending() NodeI {
	n.sortDescending = false
	return n
}

// Descending is used in an OrderBy query builder function to sort the column in descending order.
func (n *ColumnNode) Descending() NodeI {
	n.sortDescending = true
	return n
}

func (n *ColumnNode) sortDesc() bool {
	return n.sortDescending
}

/*
func (n *ColumnNode) SetValue(v interface{}) error {
	// TODO: verify
	n.value = v
	return nil
}
*/


func (n *ColumnNode) nodeType() NodeType {
	return ColumnNodeType
}

// Equals is used internally by the framework to determine if two nodes are equal.
func (n *ColumnNode) Equals(n2 NodeI) bool {
	if cn, ok := n2.(*ColumnNode); ok {
		if cn.dbTable == n.dbTable && cn.dbColumn == n.dbColumn {
			// Special code to allow new nodes to be evaluated as equal, but manual aliased nodes are not equal.
			if n.GetAlias() == "" || n2.GetAlias() == "" {
				return true
			}
			if n.GetAlias() == n2.GetAlias() {
				return true
			}
		}
	}
	return false
}

func (n *ColumnNode) name() string {
	return n.dbColumn
}

func (n *ColumnNode) goName() string {
	return n.gName
}

func (n *ColumnNode) tableName() string {
	return n.dbTable
}

func (n *ColumnNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	var alias string
	if n.alias != "" {
		alias = " as " + n.alias
	}

	log.Print(tabs + "Col: " + n.dbTable + "." + n.dbColumn + alias)
}

// GoColumnType is used internally by the framework to return the go type corresponding to the given column.
func ColumnNodeGoType(n *ColumnNode) GoColumnType {
	return n.goType
}

// ColumnNodeDbName is used internally by the framework to return the name of the column in the database.
func ColumnNodeDbName(n *ColumnNode) string {
	return n.dbColumn
}
