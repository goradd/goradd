package dbtest

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/orm/db"
)

func init() {
	//
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	//cfg.DBName = "test"
	cfg.User = "root"
	cfg.Passwd = "12345"

	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db1, key)

	db.AnalyzeDatabases()
}

