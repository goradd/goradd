package dbtest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	model2 "goradd-project/gen/goraddUnit/model"

	"testing"
)

func TestJsonMarshall1(t *testing.T) {
	ctx := getContext()
	p := model.LoadProject(ctx, "1",
		node.Project().Name(),
		node.Project().StatusType(),
		node.Project().Manager().FirstName())
	j, err := json.Marshal(p)
	assert.NoError(t, err)
	m := make(map[string]interface{})
	err = json.Unmarshal(j, &m)
	assert.NoError(t, err)
	assert.Exactly(t, "ACME Website Redesign", m["name"])
	assert.Exactly(t, "Completed", m["statusType"])
	assert.Exactly(t, "Karen", m["manager"].(map[string]interface{})["firstName"])
}

func TestJsonUnmarshall1(t *testing.T) {
	p := model.NewProject()
	err := json.Unmarshal([]byte(
		`{
	"name":"ACME Website Redesign",
	"statusType":"Completed",
	"statusTypeID":3,
	"num":14,
	"startDate":"2020-11-01T00:00:00Z"
}
`),
		&p)
	assert.NoError(t, err)
	assert.Exactly(t, "ACME Website Redesign", p.Name())
	assert.Exactly(t, model.ProjectStatusTypeCompleted, p.StatusType())
	assert.Exactly(t, 14, p.Num())
	assert.Exactly(t, 2020, p.StartDate().Year())
}

func TestJsonMarshall2(t *testing.T) {
	ctx := getContext()
	p := model.LoadPerson(ctx, "1",
		node.Person().FirstName(),
		node.Person().LastName(),
		node.Person().PersonTypes())
	j, err := json.Marshal(p)
	assert.NoError(t, err)
	m := make(map[string]interface{})
	err = json.Unmarshal(j, &m)
	assert.NoError(t, err)
	assert.Equal(t, "John", m["firstName"])
	assert.Equal(t, "Doe", m["lastName"])
	assert.ElementsMatch(t, []float64{2, 3}, m["personTypes"])
}

func TestJsonUnmarshall2(t *testing.T) {
	p := model.NewPerson()
	err := json.Unmarshal([]byte(
		`{
	"firstName":"John",
	"lastName":"Doe",
	"personTypes":[2, 3]
}
`),
		&p)
	assert.NoError(t, err)
	assert.Equal(t, "John", p.FirstName())
	assert.Equal(t, "Doe", p.LastName())
	assert.ElementsMatch(t, []model.PersonType{model.PersonTypeManager, model.PersonTypeInactive}, p.PersonTypes())
}

func TestJsonMarshall3(t *testing.T) {
	ctx := getContext()
	r := model2.LoadTypeTest(ctx, "1")
	j, err := json.Marshal(r)
	assert.NoError(t, err)

	r2 := model2.NewTypeTest()
	err = r2.UnmarshalJSON(j)
	assert.NoError(t, err)

	assert.Equal(t, 5, r2.TestInt())
	assert.Equal(t, "abcd", string(r2.TestBlob()))
}
