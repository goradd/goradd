package db

import (
	"testing"
	"goradd/gen/goradd/model"
	"goradd/gen/goradd/model/node"
	"github.com/stretchr/testify/assert"
	"context"
	. 	"github.com/spekary/goradd/orm/op"

)

func TestSubquery(t *testing.T) {
	ctx := context.Background()
	people := model.QueryPeople().
		Alias("manager_count",
		model.QueryProjects().
			Alias("", Count(node.Project().ManagerID())).
			Where(Equal(node.Project().ManagerID(), node.Person().ID())).
			Subquery()).
		Where(Equal(node.Person().LastName(), "Wolfe")).
		Load(ctx)
	assert.Equal(t, 2, people[0].GetAlias("manager_count").Int(), "Karen Wolfe manages 2 projects.")
}



