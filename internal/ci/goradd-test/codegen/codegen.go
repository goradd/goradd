package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/codegen/generator"
	_ "github.com/goradd/goradd/pkg/bootstrap/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	mysql2 "github.com/goradd/goradd/pkg/orm/db/sql/mysql"
	_ "github.com/goradd/goradd/pkg/page/control/generator"
	_ "goradd-project/config" // Initialize required variables
	_ "goradd-project/tmp/template"
)

const CiDbUser = "tester"

func main() {
	config()
	generator.Generate()
}

func config() {
	initDatabases()

}

func initDatabases() {

	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	cfg.User = CiDbUser
	cfg.Passwd = ""
	cfg.ParseTime = true

	key := "goradd"

	db1 := mysql2.NewDB(key, "", cfg)

	db.AddDatabase(db1, key)

	cfg = mysql.NewConfig()

	cfg.DBName = "goradd_unit"
	cfg.User = CiDbUser
	cfg.Passwd = ""

	key = "goraddUnit"

	db2 := mysql2.NewDB(key, "", cfg)

	db.AddDatabase(db2, key)

}
