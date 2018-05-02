package types

import (
	"testing"
)

type testObj2 struct {
	id string
}

func (t *testObj2) ID() string {
	return t.id
}

var tree *IdTree = NewIdTree()


func newObj2(id string) *testObj2 {
	c := &testObj2{id}
	return c
}

func TestBasicIdTree(t *testing.T) {
	tree := NewIdTree()

	o := tree.Get("2")
	if (o != nil) {
		t.Error("Empty tree error")
	}

	root := newObj2("root")
	tree.Add(nil, root)

	tree.Add(root, newObj2("2"))
	tree.Add(root, newObj2("3"))

	o = tree.Get("2")

	if o.ID() != "2" {
		t.Error("Could not get object from tree")
	}

	c := tree.Children(root)
	if (len(c)) != 2 {
		t.Error("Child objects not put in parent")
	}
}

func TestIdTree2(t *testing.T) {
	tree := NewIdTree()

	root := newObj2("root")
	tree.Add(nil, root)

	tree.Add(root, newObj2("b1"))
	n2 := newObj2("b2")
	tree.Add(root, n2)

	tree.Add(n2, newObj2("b2-1"))

	o := tree.Get("b2-1")

	if o.ID() != "b2-1" {
		t.Error("Could not get object from tree")
	}

	if p := tree.Parent(o); p.ID() != "b2" {
		t.Error("Parent not correct. Found: %q", p.ID())
	}

	o = tree.Root(n2)

	if o != root {
		t.Error("Could not get root")
	}

	a := tree.GetAll()
	if len(a) != 4 {
		t.Error("Could not GetAll")
	}

	tree.Remove(n2)
	if o = tree.Get("b2-1"); o != nil {
		t.Error("Could not remove branch")
	}

	if tree.Get("b2") != nil {
		t.Error("Could not remove branch")
	}

	tree.Clear()
	a = tree.GetAll()
	if len(a) != 0 {
		t.Error("Could not GetAll after a clear")
	}


}


