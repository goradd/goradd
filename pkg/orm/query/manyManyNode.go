package query

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

// A ManyManyNode is an association node that links one table to another table with a many-to-many relationship.
// Some of the columns have overloaded meanings depending on SQL or NoSQL mode.
type ManyManyNode struct {
	nodeAlias
	nodeCondition
	nodeLink
	// Which database in the global list of databases does the node belong to
	dbKey string
	// The association table
	dbTable string
	// The column in the association table pointing toward the primary object.
	dbColumn string
	// Property in the primary object used to refer to the object collection.
	goPropName string

	// The table we are joining to.
	refTable string
	// Column in association table pointing forwards to refTable
	refColumn string
	// Primary key column of refTable
	refPk string
	// Are we expanding as an array, or one item at a time.
	isArray bool
	// Is this pointing to a type table item?
	isTypeTable bool
}

// NewManyManyNode  is used internally by the ORM to create a new many-to-many node.
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
	// The primary key of refTableName
	refPk string,
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
		refPk:       refPk,
		isArray:     true,
		isTypeTable: isType,
	}
	return n
}

func (n *ManyManyNode) copy() NodeI {
	ret := &ManyManyNode{
		dbKey:         n.dbKey,
		dbTable:       n.dbTable,
		dbColumn:      n.dbColumn,
		goPropName:    n.goPropName,
		refTable:      n.refTable,
		refColumn:     n.refColumn,
		refPk:         n.refPk,
		isArray:       n.isArray,
		isTypeTable:   n.isTypeTable,
		nodeAlias:     nodeAlias{n.alias},
		nodeCondition: nodeCondition{n.condition}, // shouldn't need to duplicate condition
	}
	return ret
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

func (n *ManyManyNode) isExpander() bool {
	return true
}

// Equals is used internally by the framework to test if the node is the same as another node.
func (n *ManyManyNode) Equals(n2 NodeI) bool {
	if tn, ok := n2.(TableNodeI); !ok {
		return false
	} else if cn, ok2 := tn.EmbeddedNode_().(*ManyManyNode); !ok2 {
		return false
	} else {
		return cn.dbTable == n.dbTable &&
			cn.goPropName == n.goPropName &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)

	}
}

func (n *ManyManyNode) tableName() string {
	return n.refTable
}

func (n *ManyManyNode) databaseKey() string {
	return n.dbKey
}

func (n *ManyManyNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "MM: " + n.dbTable + "." + n.dbColumn + "." + n.refTable + "." + n.refColumn + " AS " + n.GetAlias())
}

// Return the name as a captialized object name
func (n *ManyManyNode) goName() string {
	return n.goPropName
}

type manyManyNodeEncoded struct {
	Alias       string
	Condition   NodeI
	Parent      NodeI
	DbKey       string
	DbTable     string
	DbColumn    string
	GoPropName  string
	RefTable    string
	RefColumn   string
	RefPk       string
	IsArray     bool
	IsTypeTable bool
}

func (n *ManyManyNode) GobEncode() (data []byte, err error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	s := manyManyNodeEncoded{
		Alias:       n.alias,
		Condition:   n.condition,
		Parent:      n.parentNode,
		DbKey:       n.dbKey,
		DbTable:     n.dbTable,
		DbColumn:    n.dbColumn,
		GoPropName:  n.goPropName,
		RefTable:    n.refTable,
		RefColumn:   n.refColumn,
		RefPk:       n.refPk,
		IsArray:     n.isArray,
		IsTypeTable: n.isTypeTable,
	}

	if err = e.Encode(s); err != nil {
		panic(err)
	}

	data = buf.Bytes()
	return
}

func (n *ManyManyNode) GobDecode(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var s manyManyNodeEncoded
	if err = dec.Decode(&s); err != nil {
		panic(err)
	}
	n.alias = s.Alias
	n.condition = s.Condition
	n.dbKey = s.DbKey
	n.dbTable = s.DbTable
	n.dbColumn = s.DbColumn
	n.goPropName = s.GoPropName
	n.refTable = s.RefTable
	n.refColumn = s.RefColumn
	n.refPk = s.RefPk
	n.isArray = s.IsArray
	n.isTypeTable = s.IsTypeTable
	SetParentNode(n, s.Parent)
	return
}

func init() {
	gob.Register(&ManyManyNode{})
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

// ManyManyNodeRefPk is used internally by the ORM to return the primary key column name of the table being pointed to.
func ManyManyNodeRefPk(n *ManyManyNode) string {
	return n.refPk
}

// ManyManyNodeDbTable is used internally by the framework to return the table name of the table the node belongs to
func ManyManyNodeDbTable(n *ManyManyNode) string {
	return n.dbTable
}

// ManyManyNodeDbColumn is used internally by the framework to return the column name in the table the node belongs to
func ManyManyNodeDbColumn(n *ManyManyNode) string {
	return n.dbColumn
}
