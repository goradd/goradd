package query

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

type ColumnNodeI interface {
	NodeI
	NodeSorter
	Aliaser
	nodeLinkI
}

// A Column represents a table or field in a database structure, and is the leaf of a node tree or chain.
type ColumnNode struct {
	nodeAlias
	nodeLink
	// Which database in the global list of databases does the node belong to
	dbKey string
	// Name of table in the database we point to
	dbTable string
	// The name of the column in the database
	dbColumn string
	// The name of the function used to access the property as a node or ORM item
	gName string
	// The go type for the column
	goType GoColumnType
	// Used by OrderBy clauses
	sortDescending bool
	// True if this is the private key of its parent table
	isPK bool
}

// NewColumnNode is used by the code generator to create a new column node.
func NewColumnNode(dbKey string, dbTable string, dbName string, goName string, goType GoColumnType, isPK bool) *ColumnNode {
	n := &ColumnNode{
		dbKey:    dbKey,
		dbTable:  dbTable,
		dbColumn: dbName,
		gName:    goName,
		goType:   goType,
		isPK:     isPK,
	}
	return n
}

// Returns a copy of the node, satisfying the copy interface
func (n *ColumnNode) copy() NodeI {
	ret := &ColumnNode{
		dbKey:     n.dbKey,
		dbTable:   n.dbTable,
		dbColumn:  n.dbColumn,
		gName:     n.gName,
		goType:    n.goType,
		isPK:      n.isPK,
		nodeAlias: nodeAlias{n.alias},
		// don't copy links!
	}
	return ret
}

func (n *ColumnNode) nodeType() NodeType {
	return ColumnNodeType
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

// Equals is used internally by the framework to determine if two nodes are equal.
func (n *ColumnNode) Equals(n2 NodeI) bool {
	if cn, ok := n2.(*ColumnNode); ok {
		if cn.dbTable == n.dbTable && cn.dbColumn == n.dbColumn {
			// Allow new nodes to be evaluated as equal, but manual aliased nodes are not equal.
			if n.alias == "" || cn.alias == "" {
				return true
			}
			return n.alias == cn.alias
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

func (n *ColumnNode) databaseKey() string {
	return n.dbKey
}

func (n *ColumnNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	var alias string
	if n.alias != "" {
		alias = " as " + n.alias
	}

	log.Print(tabs + "Col: " + n.dbTable + "." + n.dbColumn + alias)
}

// ColumnNodeGoType is used internally by the framework to return the go type corresponding to the given column.
func ColumnNodeGoType(n *ColumnNode) GoColumnType {
	return n.goType
}

// ColumnNodeDbName is used internally by the framework to return the name of the column in the database.
func ColumnNodeDbName(n *ColumnNode) string {
	return n.dbColumn
}

func ColumnNodeIsPK(n *ColumnNode) bool {
	return n.isPK
}

func NodeIsPK(n NodeI) bool {
	if cn, ok := n.(*ColumnNode); !ok {
		return false
	} else {
		return cn.isPK
	}
}

func (n *ColumnNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	if err = e.Encode(n.alias); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbKey); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbTable); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbColumn); err != nil {
		panic(err)
	}
	if err = e.Encode(n.gName); err != nil {
		panic(err)
	}
	if err = e.Encode(n.goType); err != nil {
		panic(err)
	}
	if err = e.Encode(n.sortDescending); err != nil {
		panic(err)
	}
	if err = e.Encode(n.isPK); err != nil {
		panic(err)
	}

	var n2 NodeI = n.nodeLink.parentNode
	err = e.Encode(&n2)
	data = buf.Bytes()
	return
}

func (n *ColumnNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&n.alias); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbKey); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbTable); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbColumn); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.gName); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.goType); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.sortDescending); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.isPK); err != nil {
		panic(err)
	}

	var n2 NodeI

	if err = dec.Decode(&n2); err != nil {
		panic(err)
	}
	SetParentNode(n, n2)
	return
}

func init() {
	gob.Register(&ColumnNode{})
}
