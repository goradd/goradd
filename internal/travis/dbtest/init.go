package dbtest

import (
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/orm/db"
)

func init() {
	cfg := mysql.NewConfig()

	cfg.DBName = "goradd"
	cfg.User = "travis"
	cfg.Passwd = ""

	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)
	//db1 := db.NewJsonLink(key, "")

	db.AddDatabase(db1, key)
}

func getContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, goradd.SqlContext, &db.SqlContext{})
	return ctx
}