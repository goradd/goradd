package dbtest

import (
	"bytes"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
)

func TestNodeEquality(t *testing.T) {

	n := node.Person()
	if !n.Equals(n) {
		t.Error("Table node not equal to self")
	}

	n = node.Project().Manager()
	if !n.Equals(n) {
		t.Error("Reference node not equal to self")
	}

	n2 := node.Person().ProjectsAsManager()
	if !n2.Equals(n2) {
		t.Error("Reverse Reference node not equal to self")
	}

	n3 := node.Person().Projects()
	if !n3.Equals(n3) {
		t.Error("Many-Many node not equal to self")
	}

	n4 := query.NewValueNode(model.PersonTypeContractor)
	if !n4.Equals(n4) {
		t.Error("Type node not equal to self")
	}

}

func BenchmarkNodeType1(b *testing.B) {
	n := node.Project().Manager()

	for i := 0; i < b.N; i++ {
		t := query.NodeGetType(n)
		if t == query.ReferenceNodeType {
			_ = n
		}
	}
}

func BenchmarkNodeType2(b *testing.B) {
	n := node.Project().Manager()

	for i := 0; i < b.N; i++ {
		if r, ok := n.EmbeddedNode_().(*query.ReferenceNode); ok {
			_ = r
		}
	}
}

func TestNodeSerialize(t *testing.T) {
	var n query.NodeI = node.Person().FirstName()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(&n)
	assert.NoError(t, err)

	var n2 query.NodeI
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&n2)
	assert.NoError(t, err)

	assert.True(t, n2.Equals(n))
}
