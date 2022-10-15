package query

/*
type FieldType int

const (
	Blob FieldType = iota
	IntegerTextbox
	VarChar
	Char
	DateTime
	Date
	Time
	FloatTextbox
	Bit
)*/

// NodeType indicates the type of node, which saves us from having to use reflection to determine this
type NodeType int

const (
	UnknownNodeType NodeType = iota
	TableNodeType
	ColumnNodeType
	ReferenceNodeType // forward reference from a foreign key
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

// Expander is the interface a node must satisfy to be able to be expanded upon,
// making a many-* relationship that creates multiple versions of the original object.
type Expander interface {
	Expand()
	isExpanded() bool
	isExpander() bool
}

type Aliaser interface {
	// SetAlias sets a unique name for the node as used in a database query.
	SetAlias(string)
	// GetAlias returns the alias that was used in a database query.
	GetAlias() string
}

// Nodes that can have an alias can mix this in
type nodeAlias struct {
	alias string
}

// SetAlias sets an alias which is an alternate name to use for the node in the result of a query.
// Aliases will generally be assigned during the query build process. You only need to assign a manual
// alias if
func (n *nodeAlias) SetAlias(a string) {
	n.alias = a

}

// GetAlias returns the alias name for the node.
func (n *nodeAlias) GetAlias() string {
	return n.alias
}

type conditioner interface {
	setCondition(condition NodeI)
	getCondition() NodeI
}

// Nodes that can have a condition can mix this in
type nodeCondition struct {
	condition NodeI
}

func (c *nodeCondition) setCondition(cond NodeI) {
	c.condition = cond
}

func (c *nodeCondition) getCondition() NodeI {
	return c.condition
}

// NodeI is the interface that all nodes must satisfy. A node is a representation of an object or a relationship
// between objects in a database that we use to create a query. It lets us abstract the structure of a database
// to be able to query any kind of database. Obviously, this doesn't work for all possible database structures, but
// it generally works well enough to solve most of the situations that you will come across.
type NodeI interface {
	// Equals returns true if the given node is equal to this node.
	Equals(NodeI) bool
	tableName() string
	log(level int)
	nodeType() NodeType
	databaseKey() string
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

func NodeDbKey(n NodeI) string {
	return n.databaseKey()
}

// NodeIsConditioner is used internally by the framework to determine if the node has a join condition.
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
		if c, ok := n.(conditioner); !ok {
			panic("cannot set condition on this type of node")
		} else {
			c.setCondition(condition)
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
		if tn.getParent() == nil {
			return false
		}
		if en, ok := tn.EmbeddedNode_().(Expander); ok {
			return en.isExpander()
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

// NodeGetType is an internal function used to get the node type without casting. This is particularly useful when
// dealing with TableNodeI types, because they are embedded in different concrete types by the code
// generator, so to get the specific node type involves extracting the embedded type and then doing a type select.
func NodeGetType(n NodeI) NodeType {
	return n.nodeType()
}

// Convenience method to see if a node is a table type of node, or a leaf node
func NodeIsTableNodeI(n NodeI) bool {
	t := NodeGetType(n)
	return t == ReferenceNodeType ||
		t == ReverseReferenceNodeType ||
		t == TableNodeType ||
		t == ManyManyNodeType
}

// Convenience method to see if a node is a reference type node. This is essentially table type nodes
// excluding an actual TableNode, since table nodes always start at the top level.
func NodeIsReferenceI(n NodeI) bool {
	t := NodeGetType(n)
	return t == ReferenceNodeType ||
		t == ReverseReferenceNodeType ||
		t == ManyManyNodeType
}

// Return the primary key of a node, if it has a primary key. Otherwise return nil.
func NodePrimaryKey(n NodeI) NodeI {
	if tn, ok := n.(TableNodeI); ok {
		return tn.PrimaryKeyNode()
	}
	return nil
}
