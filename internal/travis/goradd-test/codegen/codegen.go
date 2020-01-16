package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/codegen/generator"
	_ "github.com/goradd/goradd/pkg/bootstrap/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	_ "github.com/goradd/goradd/pkg/page/control/generator"
	_ "goradd-project/config" // Initialize required variables
	_ "goradd-project/tmp/template"
)

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
	cfg.User = "travis"
	cfg.Passwd = ""
	cfg.ParseTime = true

	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db1, key)

	cfg = mysql.NewConfig()

	cfg.DBName = "goraddUnit"
	cfg.User = "travis"
	cfg.Passwd = ""

	key = "goraddUnit"

	db2 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db2, key)

}
