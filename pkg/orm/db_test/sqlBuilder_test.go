package db

import (
	"context"
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

/*
func init() {
	//
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	//cfg.DBName = "test"
	cfg.User = "root"
	cfg.Passwd = "12345"

	key := "main"

	db1 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db1, key)

	db.AnalyzeDatabases()
}

*/

func TestSubquery2(t *testing.T) {
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
