package dbtest

import (
	. "github.com/goradd/goradd/pkg/orm/op"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"goradd-project/gen/goradd/model/node"
	"testing"
	//"github.com/goradd/goradd/pkg/orm/db"
	//"goradd/datetime"
	//"github.com/goradd/goradd/datetime"
	//"github.com/goradd/goradd/datetime"
)


func TestSubquery2(t *testing.T) {
	ctx := getContext()
	people := model.QueryPeople(ctx).
		Alias("manager_count",
			model.QueryProjects(ctx).
				Alias("", Count(node.Project().ManagerID())).
				Where(Equal(node.Project().ManagerID(), node.Person().ID())).
				Subquery()).
		Where(Equal(node.Person().LastName(), "Wolfe")).
		Load()
	assert.Equal(t, 2, people[0].GetAlias("manager_count").Int(), "Karen Wolfe manages 2 projects.")
}
