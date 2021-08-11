package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/orm/db"
)
const CiDbUser = "travis"

func init() {
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	cfg.User = CiDbUser
	cfg.Passwd = ""

	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db1, key)

	cfg = mysql.NewConfig()

	cfg.DBName = "goraddUnit"
	cfg.User = CiDbUser
	cfg.Passwd = ""

	key = "goraddUnit"

	db2 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db2, key)

}
