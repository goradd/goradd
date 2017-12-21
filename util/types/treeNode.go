package types

// An interface and structure together that can turn any struct into a tree structure with parent/child relationships
// To use it, simply embed the TreeNode structure into another structure and call Init. This TreeNodeI structure will
// be convertable to the parent object
type  TreeNodeI interface {
	AddChildNode(TreeNodeI)
	RemoveAllChildNodes()
	RemoveChildNode(TreeNodeI)
	ParentNode() TreeNodeI
	SetParent(TreeNodeI)
	TopNode() TreeNodeI
	ChildNodes() []TreeNodeI
}

// A kind of mixin for anything that controls child controls, which are all controls, but also the top level form or page
// Creates a parent/child tree of controls that is used by the drawing code to draw the controls
type TreeNode struct {
	self TreeNodeI
	parent   TreeNodeI
	children []TreeNodeI
}

// To correctly get to the top node, the top node must know about itself
func (n *TreeNode) Init(self TreeNodeI) {
	n.self = self
}

// addChildNode adds the given node as a child of the current node
func (n *TreeNode) AddChildNode(c TreeNodeI) {
	if n.self == nil {
		panic("Call Init before using the TreeNode to get the top node.")
	}
	if n.children == nil {
		n.children = []TreeNodeI{c}
	} else {
		n.children = append(n.children, c)
	}
	c.SetParent(n.self)
}

// The central removal function. Manages the entire remove process. Other removal functions should call here.
func (n *TreeNode) RemoveChildNode(c TreeNodeI) {
	for i,v := range n.children {
		if v == c {
			n.children = append(n.children[:i], n.children[i+1:]...) // remove found item from list
			break
		}
	}
}

func (n *TreeNode) RemoveAllChildNodes() {
	n.children = nil
}

func (n *TreeNode) ParentNode() TreeNodeI {
	return n.parent
}


func (n *TreeNode) SetParent(p TreeNodeI) {
	n.parent = p
}

// Return the top node in the hierarchy.
func (n *TreeNode) TopNode() TreeNodeI {
	if n.self == nil {
		panic("Call Init before using the TreeNode to get the top node.")
	}
	var f TreeNodeI = n.self
	for f.ParentNode() != nil {
		f = f.ParentNode()
	}
	return f
}

func (n *TreeNode) ChildNodes() []TreeNodeI {
	return n.children
}