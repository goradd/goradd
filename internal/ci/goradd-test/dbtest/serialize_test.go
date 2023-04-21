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

func TestNodeSerializeReverseReference(t *testing.T) {
	ctx := getContext()
	var n query.NodeI = node.Person().ProjectsAsManager()

	n2 := serNode(t, n)

	// can we still select a manager with the new node
	person := model.LoadPerson(ctx, "1", n2)
	assert.Len(t, person.ProjectsAsManager(), 1)
	assert.Equal(t, "3", person.ProjectsAsManager()[0].ID())
}

func TestNodeSerializeManyMany(t *testing.T) {
	ctx := getContext()
	var n query.NodeI = node.Person().ProjectsAsTeamMember()

	n2 := serNode(t, n)

	// can we still select project as team member
	person := model.LoadPerson(ctx, "1", n2)
	assert.Len(t, person.ProjectsAsTeamMember(), 2)
}

func serObject(t *testing.T, n interface{}) interface{} {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(&n)
	assert.NoError(t, err)

	var n2 interface{}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&n2)
	assert.NoError(t, err)
	return n2
}

func TestRecordSerializeComplex1(t *testing.T) {
	ctx := getContext()
	person := model.LoadPerson(ctx, "7",
		node.Person().ProjectsAsTeamMember(), // many many
		node.Person().ProjectsAsManager(),    // reverse
		node.Person().PersonTypes(),          // many many type
		node.Person().Login(),                // reverse unique
	)

	// Serialize and deserialize
	person2 := serObject(t, person).(*model.Person)
	assert.Len(t, person2.ProjectsAsTeamMember(), 2)
	assert.Len(t, person2.ProjectsAsManager(), 2)
	assert.Len(t, person2.PersonTypes(), 2)
	assert.Equal(t, "kwolfe", person2.Login().Username())
}

func TestRecordSerializeComplex2(t *testing.T) {
	ctx := getContext()
	login := model.LoadLogin(ctx, "4",
		node.Login().Person().ProjectsAsTeamMember(), // many many
		node.Login().Person().ProjectsAsManager(),    // reverse
		node.Login().Person().PersonTypes(),          // many many type
		node.Login().Person().Login(),                // reverse unique
	)

	// Serialize and deserialize
	login2 := serObject(t, login).(*model.Login)
	assert.Len(t, login2.Person().ProjectsAsTeamMember(), 2)
	assert.Len(t, login2.Person().ProjectsAsManager(), 2)
	assert.Len(t, login2.Person().PersonTypes(), 2)
	assert.Equal(t, "kwolfe", login2.Person().Login().Username())
}
