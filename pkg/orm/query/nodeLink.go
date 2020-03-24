package query

// The nodeLinkI interface provides an interface to allow nodes to be linked in a parent child chain
type nodeLinkI interface {
	setChild(NodeI)
	setParent(NodeI)
	getParent() NodeI
	getChild() NodeI
	copy() NodeI // all linkable nodes must be copyable
}

// The nodeLink is designed to be a mixin for the basic node structure. It encapsulates the joining of nodes.
// In particular, the SetParentNode method gets exported for codegen purposes.
type nodeLink struct {
	// Parent of the join, so its doubly linked
	parentNode NodeI
	// child of the join.
	childNode NodeI
}

func (n *nodeLink) setChild(cn NodeI) {
	if n.childNode == nil {
		n.childNode = cn
	} else {
		panic("node already has a child node")
	}
}

func (n *nodeLink) setParent(pn NodeI) {
	n.parentNode = pn
}

// SetParentNode is used internally by the framework.
// It is used by the codegenerator to create linked nodes.
// It is used by the serializer to restore linked nodes.
func SetParentNode(child NodeI, parent NodeI) {
	if parent != nil {
		if parent.(nodeLinkI).getChild() != nil {
			// Create a copy of the parent chain, since the parent already has a child
			parent = copyUp(parent)
		}
		child.(nodeLinkI).setParent(parent)
		parent.(nodeLinkI).setChild(child)
	}
}

// copyUp creates a copy of the given node and copies all of its parent nodes too, putting the copies in its parent chain.
func copyUp(n NodeI) NodeI {
	nl := n.(TableNodeI)
	cp := nl.Copy_()
	if p := nl.getParent(); p != nil {
		parent := copyUp(p)
		cp.(nodeLinkI).setParent(parent)
		parent.(nodeLinkI).setChild(cp)
	}
	return cp
}

// ParentNode is used internally by the framework to return a node's parent.
func ParentNode(n NodeI) NodeI {
	return n.(nodeLinkI).getParent()
}

func (n *nodeLink) getParent() NodeI {
	if n.parentNode == nil {
		return nil
	}
	return n.parentNode.(NodeI)
}

func (n *nodeLink) getChild() NodeI {
	return n.childNode
}

/**

Public Accessors

The following functions are designed primarily to be used by the db package to help it unpack queries. The are not
given an accessor at the beginning so that they do not show up as a function in editors that provide code hinting when
trying to put together a node chain during the code creation process. Essentially they are trying to create exported
functions for the db package without broadcasting them to the world.

*/

// ChildNode is used internally by the framework to get the child node of a node
func ChildNode(n NodeI) NodeI {
	return n.(nodeLinkI).getChild()
}

// RootNode is used internally by the framework to get the root node, which is the top parent in the node tree.
func RootNode(n NodeI) NodeI {
	if self, ok := n.(nodeLinkI); !ok {
		return nil
	} else if self.getParent() == nil {
		return self.(NodeI)
	} else {
		var n1 = self
		for pn := n1.getParent(); pn != nil; pn = n1.getParent() {
			n1 = pn.(nodeLinkI)
		}
		return n1.(NodeI)
	}
}

func CopyNode(n NodeI) ReferenceNodeI {
	if self, ok := n.(ReferenceNodeI); !ok {
		panic("cannot copy this kind of node")
	} else {
		return self.copy().(ReferenceNodeI)
	}
}
