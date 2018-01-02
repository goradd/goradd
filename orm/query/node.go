package query

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

type goNamer interface {
	goName() string
}

type nodeContainer interface {
	containedNodes() (nodes []NodeI)
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


/**

Public Accessors

The following functions are designed primarily to be used by the db package to help it unpack queries. The are not
given an accessor at the beginning so that they do not show up as a function in editors that provide code hinting when
trying to put together a node chain during the code creation process. Essentially they are trying to create exported
functions for the db package without broadcasting them to the world.

 */

func NodeTableName (n NodeI) string {
	return n.tableName()
}

func NodeIsConditioner(n NodeI) bool {
	if tn, _ := n.(TableNodeI); tn != nil {
		if c, _ := tn.EmbeddedNode_().(conditioner); c != nil {
			return true
		}
	}
	return false
}

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

func NodeCondition(n NodeI) NodeI {
	if tn, ok := n.(TableNodeI); ok {
		if cn, ok := tn.EmbeddedNode_().(conditioner); ok {
			return cn.getCondition()
		}
	}
	return nil
}


func ContainedNodes(n NodeI) (nodes []NodeI) {
	if nc, ok := n.(nodeContainer); ok {
		return nc.containedNodes()
	} else {
		return nil
	}
}

func LogNode(n NodeI, level int) {
	n.log(level)
}

func GetNodeType(n NodeI) NodeType {
	return n.nodeType()
}

func NodeIsExpanded(n NodeI) bool {
	if tn, ok := n.(TableNodeI); ok {
		if en,ok := tn.EmbeddedNode_().(Expander); ok {
			return en.isExpanded()
		}
	}
	return false
}

func NodeIsExpander(n NodeI) bool {
	if tn, ok := n.(TableNodeI); ok {
		if _,ok := tn.EmbeddedNode_().(Expander); ok {
			return true
		}
	}
	return false
}

func ExpandNode(n NodeI) {
	if tn, ok := n.(TableNodeI); ok {
		if en,ok := tn.EmbeddedNode_().(Expander); ok {
			en.Expand()
		} else {
			panic ("Cannot expand a node of this type")
		}
	} else {
		panic ("Cannot expand a node of this type")
	}
}

func NodeGoName(n NodeI) string {
	if gn,ok := n.(goNamer); ok {
		return gn.goName()
	} else if tn,ok := n.(TableNodeI); ok {
		return tn.EmbeddedNode_().(goNamer).goName()
	}
	return ""
}

func NodeSorterSortDesc(n NodeSorter) bool {
	return n.sortDesc()
}