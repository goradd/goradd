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

// serialize and deserialize the node
func serNode(t *testing.T, n query.NodeI) query.NodeI {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(&n)
	assert.NoError(t, err)

	var n2 query.NodeI
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&n2)
	assert.NoError(t, err)
	return n2
}

func TestNodeSerializeReference(t *testing.T) {
	ctx := getContext()
	var n query.NodeI = node.Project().Manager()

	n2 := serNode(t, n)

	// can we still select a manager with the new node
	proj := model.LoadProject(ctx, "1", n2)
	assert.Equal(t, proj.Manager().LastName(), "Wolfe")
}
