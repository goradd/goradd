package dbtest

import (
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/orm/db"
	mysql2 "github.com/goradd/goradd/pkg/orm/db/sql/mysql"
)

const CiDbUser = "tester"

func init() {
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	cfg.User = CiDbUser
	cfg.Passwd = ""

	key := "goradd"

	db1 := mysql2.NewMysqlDB(key, "", cfg)

	db.AddDatabase(db1, key)

	cfg = mysql.NewConfig()

	cfg.DBName = "goradd_unit"
	cfg.User = CiDbUser
	cfg.Passwd = ""

	key = "goraddUnit"

	db2 := mysql2.NewMysqlDB(key, "", cfg)

	db.AddDatabase(db2, key)

}

func getContext() context.Context {
	ctx := context.Background()
	for _, d := range db.GetDatabases() {
		ctx = d.PutBlankContext(ctx)
	}
	return ctx
}
