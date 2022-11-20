//go:build mysql

package dbtest

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/orm/db"
	mysql2 "github.com/goradd/goradd/pkg/orm/db/sql/mysql"
)

const CiDbUser = "tester"
const CiDbPassword = ""

func init() {
	initMysql()
}

func initMysql() {
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	cfg.User = CiDbUser
	cfg.Passwd = CiDbPassword
	key := "goradd"
	cfg.ParseTime = true

	db1 := mysql2.NewDB(key, "", cfg)
	db1.Analyze(mysql2.DefaultOptions())

	db.AddDatabase(db1, key)

	cfg = mysql.NewConfig()

	key = "goraddUnit"
	cfg.DBName = "goradd_unit"
	cfg.User = CiDbUser
	cfg.Passwd = CiDbPassword
	cfg.ParseTime = true

	db2 := mysql2.NewDB(key, "", cfg)
	db2.Analyze(mysql2.DefaultOptions())

	db.AddDatabase(db2, key)
}
