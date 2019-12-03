package query

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

// ReverseReferenceNode creates a reverse reference node representing a one to many relationship or one-to-one
// relationship, depending on whether the foreign key is unique. The other side of the relationship will have
// a matching forward ReferenceNode.
//
// For SQL databases, a ReverseReferenceNode will not have anything in its
// table to indicate the relationship. It is a kind of virtual placeholder to indicate that a foreign key in
// another table is pointing to this table, and therefore that relationship can be used to build a query.
type ReverseReferenceNode struct {
	nodeAlias
	nodeCondition
	nodeLink
	// Which database in the global list of databases does the node belong to
	dbKey string
	// dbTable is the table which owns the reverse reference
	dbTable string
	// dbKeyColumn is the primary key of that table
	dbKeyColumn string
	// dbColumn is NoSQL only. It is the table containing an array of primary keys for the records pointing back to this one.
	dbColumn string
	// goPropName is the property we are using to refer to the many side of the relationship
	goPropName string
	// refTable is the table containing the pointer back to us
	refTable string
	// refColumn is the column that is the foreign key pointing back to us.
	refColumn string
	// isArray is true to create new objects for each joined item, or false to create an array of joined objects here.
	isArray bool
}

func NewReverseReferenceNode(
	dbKey string,
	// dbTable is the name of the database table which owns the reverse reference
	dbTable string,
	// dbKeyColumn is the primary key of that table
	dbKeyColumn string,
	// dbColumn is for NoSQL and is the column containing an array of primary keys that point to the items referring back to us.
	dbColumn string,
	// Property we are using to refer to the many side of the relationship
	goName string,
	// The table containing the pointer back to us
	refTable string,
	// The column that is the foreign key pointing back to us.
	refColumn string,
	isArray bool,
) *ReverseReferenceNode {
	n := &ReverseReferenceNode{
		dbKey:      dbKey,
		dbTable:    dbTable,
		dbKeyColumn: dbKeyColumn,
		dbColumn:   dbColumn,
		goPropName: goName,
		refTable:   refTable,
		refColumn:  refColumn,
		isArray:    isArray,
	}
	return n
}

func (n *ReverseReferenceNode) copy() NodeI {
	ret := &ReverseReferenceNode{
		dbKey:         n.dbKey,
		dbTable:       n.dbTable,
		dbKeyColumn:   n.dbKeyColumn,
		dbColumn:      n.dbColumn,
		goPropName:    n.goPropName,
		refTable:      n.refTable,
		refColumn:     n.refColumn,
		isArray:       n.isArray,
		nodeAlias:     nodeAlias{n.alias},
		nodeCondition: nodeCondition{n.condition},
	}
	return ret
}

// Expand tells the node to expand its results into multiple records, one for each item found in this relationship,
// rather than have this relationship create an array of items within an individual record.
// Unique reverse relationships create one-to-one relationships, and so they are always expanded.
func (n *ReverseReferenceNode) Expand() {
	n.isArray = false
}

func (n *ReverseReferenceNode) isExpanded() bool {
	return !n.isArray
}

// Equals is used internally by the framework to determine if two nodes are equal.
func (n *ReverseReferenceNode) Equals(n2 NodeI) bool {
	if tn, ok := n2.(TableNodeI); !ok {
		return false
	} else if cn, ok := tn.EmbeddedNode_().(*ReverseReferenceNode); !ok {
		return false
	} else {
		return cn.dbTable == n.dbTable &&
			cn.goPropName == n.goPropName &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)
	}
}

func (n *ReverseReferenceNode) tableName() string {
	return n.refTable
}

func (n *ReverseReferenceNode) databaseKey() string {
	return n.dbKey
}


func (n *ReverseReferenceNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "RR: " + n.dbTable + "." + n.refTable + "." + n.refColumn + " AS " + n.GetAlias())
}

// Return the name as a captialized object name
func (n *ReverseReferenceNode) goName() string {
	return n.goPropName
}

func (n *ReverseReferenceNode) nodeType() NodeType {
	return ReverseReferenceNodeType
}

func (n *ReverseReferenceNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	if err = e.Encode(n.alias); err != nil {
		panic(err)
	}
	if err = e.Encode(n.condition); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbKey); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbTable); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbKeyColumn); err != nil {
		panic(err)
	}
	if err = e.Encode(n.dbColumn); err != nil {
		panic(err)
	}
	if err = e.Encode(n.goPropName); err != nil {
		panic(err)
	}
	if err = e.Encode(n.refTable); err != nil {
		panic(err)
	}
	if err = e.Encode(n.refColumn); err != nil {
		panic(err)
	}
	if err = e.Encode(n.isArray); err != nil {
		panic(err)
	}

	err = e.Encode(n.nodeLink.parentNode)
	data = buf.Bytes()
	return
}


func (n *ReverseReferenceNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&n.alias); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.condition); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbKey); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbTable); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbKeyColumn); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.dbColumn); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.goPropName); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.refTable); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.refColumn); err != nil {
		panic(err)
	}
	if err = dec.Decode(&n.isArray); err != nil {
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
	gob.Register(&ReverseReferenceNode{})
}


// ReverseReferenceNodeIsArray is used internally by the framework to determine if a node should create an array
func ReverseReferenceNodeIsArray(n *ReverseReferenceNode) bool {
	return n.isArray
}

// ReverseReferenceNodeRefTable is used internally by the framework to get the referenced table
func ReverseReferenceNodeRefTable(n *ReverseReferenceNode) string {
	return n.refTable
}

// ReverseReferenceNodeRefColumn is used internally by the framework to get the referenced column
func ReverseReferenceNodeRefColumn(n *ReverseReferenceNode) string {
	return n.refColumn
}

// ReverseReferenceNodeDbColumnName is used internally by the framework to get the database column on this side of the
// relationship. This is for NoSQL only
func ReverseReferenceNodeDbColumnName(n *ReverseReferenceNode) string {
	return n.dbColumn
}

// ReverseReferenceNodeKeyColumnName is used internally by the framework to get the database column on this side of the
// relationship, which is most likely the primary key. This is for SQL databases generally.
func ReverseReferenceNodeKeyColumnName(n *ReverseReferenceNode) string {
	return n.dbKeyColumn
}

