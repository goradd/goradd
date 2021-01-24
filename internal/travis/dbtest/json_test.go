package dbtest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
)

func TestJsonMarshall1(t *testing.T) {
	ctx := getContext()
	p := model.LoadProject(ctx, "1",
		node.Project().Name(),
		node.Project().ProjectStatusType(),
		node.Project().Manager().FirstName())
	j,err := json.Marshal(p)
	assert.NoError(t, err)
	m := make(map[string]interface{})
	err = json.Unmarshal(j, &m)
	assert.NoError(t, err)
	assert.Equal(t, "ACME Website Redesign", m["name"])
	assert.Equal(t, "3", m["projectStatusType"])
	assert.Equal(t, "Karen", m["manager"].(map[string]interface{})["firstName"])
}

func TestJsonUnmarshall1(t *testing.T) {
	p := model.NewProject()
	err := json.Unmarshal([]byte(
		`{
	"name":"ACME Website Redesign",
	"projectStatusType":"3",
	"num":14,
	"startDate":"2020-11-01"
}
`),
		&p)
	assert.NoError(t, err)
	assert.Equal(t, "ACME Website Redesign", p.Name())
	assert.Equal(t, model.ProjectStatusTypeCompleted, p.ProjectStatusType())
	assert.Equal(t, 14, p.Num())
	assert.Equal(t, 2020, p.StartDate().Year())
}

func TestJsonMarshall2(t *testing.T) {
	ctx := getContext()
	p := model.LoadPerson(ctx, "1",
		node.Person().FirstName(),
		node.Person().LastName(),
		node.Person().PersonTypes())
	j,err := json.Marshal(p)
	assert.NoError(t, err)
	m := make(map[string]interface{})
	err = json.Unmarshal(j, &m)
	assert.NoError(t, err)
	assert.Equal(t, "John", m["firstName"])
	assert.Equal(t, "Doe", m["lastName"])
	assert.ElementsMatch(t, []string{"2","3"}, m["personTypes"])
}

func TestJsonUnmarshall2(t *testing.T) {
	p := model.NewPerson()
	err := json.Unmarshal([]byte(
		`{
	"firstName":"John",
	"lastName":"Doe",
	"personTypes":["2", "3"]
}
`),
		&p)
	assert.NoError(t, err)
	assert.Equal(t, "John", p.FirstName())
	assert.Equal(t, "Doe", p.LastName())
	assert.ElementsMatch(t, []model.PersonType{model.PersonTypeManager, model.PersonTypeInactive}, p.PersonTypes())
}

