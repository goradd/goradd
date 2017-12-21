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

)

type nodeContainer interface {
	containedNodes() []NodeI
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
	nodeType() NodeType
	setAlias(string)
	getAlias() string
	tableName() string
	log(level int)
}

type Node struct {
	nodeLink
	conditions []NodeI	// Used only by expansion nodes
	alias string
}

type conditioner interface {
	setConditions(conditions []NodeI)
	getConditions() []NodeI
}

func (n *Node) setAlias(a string) {
	n.alias = a
}

func (n *Node) getAlias() string {
	return n.alias
}
