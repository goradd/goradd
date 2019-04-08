package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/orm/db"
)

func init() {
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	cfg.User = "travis"
	cfg.Passwd = ""

	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db1, key)
}

