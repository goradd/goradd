package db
import (
)

type FieldType int

const (
	Blob FieldType = iota
	Integer
	VarChar
	Char
	DateTime
	Date
	Time
	Float
	Bit
)

type NodeType int

const (
	UNKNOWN_NODE_TYPE		NodeType = iota
	TABLE_NODE
	COLUMN_NODE
	REFERENCE_NODE                                                     // forward reference from a foreign key
	MANYMANY_NODE
	REVERSE_REFERENCE_NODE
	VALUE_NODE
	OPERATION_NODE
	ALIAS_NODE
	SUBQUERY_NODE

)

type nodeContainer interface {
	containedNodes() (nodes []NodeI)
}

type objectNode interface {
	objectName() string
}

type NodeSorter interface {
	Ascending() NodeI
	Descending() NodeI
	sortDesc() bool
}

type Expander interface {
	Expand()
	isExpanded() bool
}

type NodeI interface {
	nodeLinkI
	Equals(NodeI) bool
	SetAlias(string)
	GetAlias() string
	nodeType() NodeType
	tableName() string
	log(level int)
}

type Node struct {
	nodeLink
	condition NodeI	// Used only by expansion nodes
	alias string
}

type conditioner interface {
	setCondition(condition NodeI)
	getCondition() NodeI
}

func (n *Node) SetAlias(a string) {
	n.alias = a
}

func (n *Node) GetAlias() string {
	return n.alias
}
