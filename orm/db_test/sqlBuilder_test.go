package db

import (
	"testing"
	"grlocal/model"
	"github.com/stretchr/testify/assert"
	"context"
	. 	"github.com/spekary/goradd/orm/op"
	"grlocal/model/node"
	//"goradd/orm/db"
	//"goradd/datetime"
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



func TestDeleteQuery(t *testing.T) {
	ctx := context.Background()

	person := model.NewPerson()
	person.SetFirstName("Test1")
	person.SetLastName("Last1")
	person.Save(ctx)

	model.QueryPeople().
		Where(
		And(
			Equal(
				node.Person().FirstName(), "Test1"),
			Equal(
				node.Person().LastName(), "Last1"))).
		Delete(ctx)

	people := model.QueryPeople().
		Where(
		And(
			Equal(
				node.Person().FirstName(), "Test1"),
			Equal(
				node.Person().LastName(), "Last1"))).
		Load(ctx)

	assert.Len(t, people, 0, "Deleted the person")
}
