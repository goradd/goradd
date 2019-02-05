package query

// The nodeLinkI interface provides an interface to allow nodes to be linked in a parent multi-child chain
type nodeLinkI interface {
	addChildNode(NodeI)
	setSelf(NodeI)
	setParent(NodeI)
	getParentNode() NodeI
	rootNode() NodeI
	getChildNodes() []NodeI
}

// The nodeLink is designed to be a mixin for the basic node structure. It encapsulates the chaining of nodes.
// In particular, the SetParentNode method gets exported for codegen purposes
type nodeLink struct {
	// We need to be able to get to the surrounding class. This is our aid to doing that. It must be set up correctly though.
	self NodeI
	// Parent in the chain, so its doubly linked
	parentNode NodeI
	// Child nodes. Multiple child nodes are used to indicate where in the chain of a query we are expanding.
	childNodes []NodeI
}

func (n *nodeLink) addChildNode(cn NodeI) {
	if n.childNodes == nil {
		n.childNodes = []NodeI{cn}
	} else {
		var found = false
		for _, n2 := range n.childNodes {
			if n2.Equals(cn) {
				found = true
			}
		}
		if !found {
			n.childNodes = append(n.childNodes, cn)
		}
	}
}

func (n *nodeLink) setSelf(self NodeI) {
	if n.self == nil {
		n.self = self
	}
}

func (n *nodeLink) setParent(p NodeI) {
	n.parentNode = p
}

// SetParentNode is used internally by the framework.
// It is used by the query builder to build a join tree, and by the codegenerator to initialize nodes
func SetParentNode(child NodeI, parent NodeI) {
	if parent != nil {
		child.setParent(parent)
		parent.addChildNode(child)
	}
	child.setSelf(child)
}

// ParentNode is used internally by the framework to return a node's parent.
func ParentNode(n NodeI) NodeI {
	return n.getParentNode()
}

// rootNode returns the top node in the node chain
func (n *nodeLink) rootNode() NodeI {
	if n.self == nil {
		return nil // value or operation node
	}
	var n1 NodeI = n.self
	for pn := n1.getParentNode(); pn != nil; pn = n1.getParentNode() {
		n1 = pn
	}
	return n1
}

func (n *nodeLink) getParentNode() NodeI {
	return n.parentNode
}

func (n *nodeLink) getChildNodes() []NodeI {
	return n.childNodes
}

/**

Public Accessors

The following functions are designed primarily to be used by the db package to help it unpack queries. The are not
given an accessor at the beginning so that they do not show up as a function in editors that provide code hinting when
trying to put together a node chain during the code creation process. Essentially they are trying to create exported
functions for the db package without broadcasting them to the world.

*/

// ChildNodes is used internally by the framework to get the child nodes of a ndoe
func ChildNodes(n nodeLinkI) []NodeI {
	return n.getChildNodes()
}

// RootNode is used internally by the framework to get the root node, which is the top parent in the node tree.
func RootNode(n nodeLinkI) NodeI {
	return n.rootNode()
}
