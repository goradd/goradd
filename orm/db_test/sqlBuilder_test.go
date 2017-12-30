package db
/*
import (
	"testing"
	"grlocal/model"
	"github.com/stretchr/testify/assert"
	"context"
	//. 	"github.com/spekary/goradd/orm/op"
	"grlocal/model/node"
	//"github.com/spekary/goradd/orm/db"
	//"goradd/datetime"
	//"github.com/spekary/goradd/datetime"
	//"github.com/spekary/goradd/datetime"
)

*/

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

/*
func TestSubquery(t *testing.T) {
	ctx := context.Background()
	people := model.QueryPeople().
		Alias("manager_count",
			Count(false,
				Subquery(model.QueryProjects().
					Where(Equal(node.Project().ManagerID(), node.Person().ID()))))).
		Where(Equal(node.Person().LastName(), "Wolfe")).
		Load(ctx)
	assert.Equal(t, 2, people[0].GetAlias("manager_count").Int(), "Karen Wolfe manages 2 projects.")
}


*/
