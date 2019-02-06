package query

/*
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
)*/


// NodeType indicates the type of node, which saves us from having to use reflection to determine this
type NodeType int

const (
	UnknownNodeType NodeType = iota
	TableNodeType
	ColumnNodeType
	ReferenceNodeType  // forward reference from a foreign key
	ManyManyNodeType
	ReverseReferenceNodeType
	ValueNodeType
	OperationNodeType
	AliasNodeType
	SubqueryNodeType
)

type goNamer interface {
	goName() string
}

type nodeContainer interface {
	containedNodes() (nodes []NodeI)
}

// NodeSorter is the interface a node must satisfy to be able to be used in an OrderBy statement.
type NodeSorter interface {
	Ascending() NodeI
	Descending() NodeI
	sortDesc() bool
}

// Expander is the interface a node must satisfy to be able to be expanded upon, making a many-* relationship create multiple versions of the original object.
type Expander interface {
	Expand()
	isExpanded() bool
}

// NodeI is the interface that all nodes must satisfy
type NodeI interface {
	nodeLinkI
	// Equals returns true if the given node is equal to this node
	Equals(NodeI) bool
	// SetAlias sets a unique name for the node as used in a database query
	SetAlias(string)
	// GetAlias returns the alias that was used in a database query
	GetAlias() string
	nodeType() NodeType
	tableName() string
	log(level int)
}

// Node is the base mixin for all node structures. A node is a representation of an object or a relationship
// between objects in a database that we use to create a query. It lets us abstract the structure of a database
// to be able to query any kind of database. Obviously, this doesn't work for all possible database structures, but
// it generally works well enough to solve most, if not all, of the situations that you will come across.
type Node struct {
	nodeLink
	condition NodeI // Used only by expansion nodes
	alias     string
}

type conditioner interface {
	setCondition(condition NodeI)
	getCondition() NodeI
}

// SetAlias sets an alias which is an alternate name to use for the node in the result of a query.
func (n *Node) SetAlias(a string) {
	n.alias = a
}

// GetAlias returns the alias name for the node.
func (n *Node) GetAlias() string {
	return n.alias
}

/**

Public Accessors

The following functions are designed primarily to be used by the db package to help it unpack queries. They are not
given an accessor at the beginning so that they do not show up as a function in editors that provide code hinting when
trying to put together a node chain during the code creation process. Essentially they are trying to create exported
functions for the db package without broadcasting them to the world.

*/

// NodeTableName is used internally by the framework to return the table associated with a node.
func NodeTableName(n NodeI) string {
	return n.tableName()
}

// NodeIsConditioner is used internally by the framework to determine if the node has a condition.
func NodeIsConditioner(n NodeI) bool {
	if tn, _ := n.(TableNodeI); tn != nil {
		if c, _ := tn.EmbeddedNode_().(conditioner); c != nil {
			return true
		}
	}
	return false
}

// NodeSetCondition is used internally by the framework to set a condition on a node.
func NodeSetCondition(n NodeI, condition NodeI) {
	if condition != nil {
		if tn, ok := n.(TableNodeI); ok {
			if c, ok := tn.EmbeddedNode_().(conditioner); !ok {
				panic("Cannot set condition on this type of node")
			} else {
				c.setCondition(condition)
			}
		} else {
			panic("Cannot set condition on this type of node")
		}
	}
}

// NodeCondition is used internally by the framework to get a condition node.
func NodeCondition(n NodeI) NodeI {
	if cn, ok := n.(conditioner); ok {
		return cn.getCondition()
	} else if tn, ok := n.(TableNodeI); ok {
		if cn, ok := tn.EmbeddedNode_().(conditioner); ok {
			return cn.getCondition()
		}
	}
	return nil
}

// ContainedNodes is used internally by the framework to return the contained nodes.
func ContainedNodes(n NodeI) (nodes []NodeI) {
	if nc, ok := n.(nodeContainer); ok {
		return nc.containedNodes()
	} else {
		return nil
	}
}

// LogNode is used internally by the framework to debug node issues.
func LogNode(n NodeI, level int) {
	n.log(level)
}

// GetNodeType is used internally by the framework to get the type of node without using reflection.
func GetNodeType(n NodeI) NodeType {
	return n.nodeType()
}

// NodeIsExpanded is used internally by the framework to detect if the node is an expanded many-many relationship.
func NodeIsExpanded(n NodeI) bool {
	if tn, ok := n.(TableNodeI); ok {
		if en, ok := tn.EmbeddedNode_().(Expander); ok {
			return en.isExpanded()
		}
	}
	return false
}

// NodeIsExpander is used internally by the framework to detect if the node can be an expanded many-many relationship.
func NodeIsExpander(n NodeI) bool {
	if tn, ok := n.(TableNodeI); ok {
		if _, ok := tn.EmbeddedNode_().(Expander); ok {
			return true
		}
	}
	return false
}

// ExpandNode is used internally by the framework to expand a many-many relationship.
func ExpandNode(n NodeI) {
	if tn, ok := n.(TableNodeI); ok {
		if en, ok := tn.EmbeddedNode_().(Expander); ok {
			en.Expand()
		} else {
			panic("Cannot expand a node of this type")
		}
	} else {
		panic("Cannot expand a node of this type")
	}
}

// NodeGoName is used internally by the framework to return the go name of the item the node refers to.
func NodeGoName(n NodeI) string {
	if gn, ok := n.(goNamer); ok {
		return gn.goName()
	} else if tn, ok := n.(TableNodeI); ok {
		return tn.EmbeddedNode_().(goNamer).goName()
	}
	return ""
}

// NodeSorterSortDesc is used internally by the framework to determine if the NodeSorter is descending.
func NodeSorterSortDesc(n NodeSorter) bool {
	return n.sortDesc()
}
