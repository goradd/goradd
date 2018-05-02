package types
import (
	"testing"
	"strconv"
)

type testObj struct {
	TreeNode
	id string
}

type testObjI interface {
	getID() string
}

var idCounter int = 0

func newObj() *testObj {
	c := &testObj{}
	c.id = strconv.Itoa(idCounter)
	idCounter++
	c.Init(c)
	return c
}

func (o *testObj) getID() string {
	return o.id
}

func TestTreeNodeBasic (t *testing.T) {
	top := newObj()

	top.AddChildNode(newObj())
	top.AddChildNode(newObj())
	top.AddChildNode(newObj())

	children := top.ChildNodes()

	if len(children) != 3 {
		t.Error("Wrong number of children.")
	}

	if _,ok := children[0].(*testObj); !ok {
		t.Error("Could not case object to parent.")
	}

	if _,ok := children[0].(testObjI); !ok {
		t.Error("Could not case object to parent interface.")
	}

	children[0].AddChildNode(newObj())

	children2 := children[0].ChildNodes()

	if top := children2[0].TopNode(); top == nil {
		t.Error("Could not get top node.")
	}

	if i,ok := children2[0].TopNode().(testObjI); !ok {
		t.Error("Could not cast top node.")
	} else if i.getID() != "0" {
		t.Error("Top node is not correct.")
	}
}