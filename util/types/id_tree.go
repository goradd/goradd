package types

import (
	"fmt"
	"sync"
)

// The IdTree is a collection of Ider objects in a tree structure. Objects must have a unique id within the structure.
// It is likely that each of your objects will need a pointer to the IdTree to manipulate it, but that is an implementation
// dependant thing
type IdTree struct {
	nodes map[string]node
	sync.RWMutex
}

// Any object that returns a string id can be stored in the tree.
type Ider interface {
	Id() string
}

type node struct {
	parentId string
	children []string
	value Ider
}

func NewIdTree() *IdTree {
	return &IdTree{nodes: make (map[string]node)}
}

func (t *IdTree) Get(id string) Ider {
	t.Lock()
	defer t.Unlock()

	if node,ok := t.nodes[id]; ok {
		return node.value
	} else {
		return nil
	}
}


// Add an item to the tree.
func (t *IdTree) Add(parent Ider, child Ider) {
	var childId string

	if childId = child.Id(); childId == "" {
		panic ("An item must have an id before it can be added to the tree")
	}

	t.Lock()
	defer t.Unlock()

	if _,ok := t.nodes[childId]; ok {
		panic(fmt.Sprintf("An item already exists with the id: %q", childId))
	}

	var parentId string

	if parent != nil {
		if parentId = parent.Id(); parentId == "" {
			panic ("Parent must have an id.")
		}

		if parentNode,ok := t.nodes[parentId]; !ok {
			panic("Parent not found in tree.")
		} else {
			n := node{parentId, nil, child}
			t.nodes[childId] = n;
			if parentNode.children == nil {
				parentNode.children = []string{childId}
			} else {
				parentNode.children = append(parentNode.children, childId)
			}
			t.nodes[parentNode.value.Id()] = parentNode
		}
	} else {
		// a top level item
		n := node{"", nil, child}
		t.nodes[childId] = n;
	}
}


// Removes the item from the tree and all its sub-itmes
func (t *IdTree) Remove(item Ider) {
	var id string = item.Id()
	if id == "" {
		panic("The item to remove does not have an id")
	}

	t.Lock()
	defer t.Unlock()

	t.remove(id)
}

// Remove all child controls
func (t *IdTree) RemoveChildren(parent Ider) {
	var id string = parent.Id()
	if id == "" {
		panic("The item does not have an id")
	}
	t.Lock()
	defer t.Unlock()

	if node,ok := t.nodes[id]; !ok {
		panic("The item is not in the tree.")
	} else if node.children != nil {
		for _,id2 := range node.children {
			t.remove(id2)
		}
	}
}


func (t *IdTree) remove(id string) {
	var n node
	var ok bool

	if n, ok = t.nodes[id]; !ok {
		panic("The item to remove was not found")
	}

	if n.children != nil {
		for _, childId := range n.children {
			t.remove(childId) // recurse
		}
	}
	delete (t.nodes, id)
}

// GetAll returns a list of all the Ider items in the tree. Order is random.
func (t *IdTree) GetAll() []Ider {
	l := make([]Ider, len(t.nodes))
	i := 0

	t.Lock()
	defer t.Unlock()

	for _,n := range t.nodes {
		l[i] = n.value
		i++
	}
	return l
}

func (t *IdTree) Children(parent Ider) []Ider {
	var parentId string = parent.Id()
	var n node
	var ok bool

	if n, ok = t.nodes[parentId]; !ok {
		panic("The item was not found")
	}

	if n.children == nil {
		return []Ider{}
	}

	l := make([]Ider, len(n.children))

	t.Lock()
	defer t.Unlock()

	i := 0
	for _,childId := range n.children {
		l[i] = t.nodes[childId].value
		i++
	}

	return l
}

func (t *IdTree) Parent(child Ider) Ider {
	var childId string = child.Id()
	var ok bool
	var n node

	if n, ok = t.nodes[childId]; !ok {
		panic("The item was not found")
	}

	t.Lock()
	defer t.Unlock()
	return t.nodes[n.parentId].value
}

// Root returns the root of the branch that the given object is on
func (t *IdTree) Root(child Ider) Ider {
	var childId string = child.Id()
	var ok bool
	var n node

	if n, ok = t.nodes[childId]; !ok {
		panic("The item was not found")
	}

	t.Lock()
	defer t.Unlock()
	for n.parentId != "" {
		n = t.nodes[n.parentId];
	}
	return n.value
}

func (t *IdTree) Clear() {
	t.Lock()
	defer t.Unlock()

	t.nodes = make (map[string]node)
}
