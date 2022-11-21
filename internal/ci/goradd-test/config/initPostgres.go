package config

import (
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/db/sql/pgsql"
)

const CiDbUser = "root"
const CiDbPassword = "12345"

func initDatabases() {
	cfg, _ := pgx.ParseConfig("")

	key := "goradd"
	cfg.Host = "localhost"
	cfg.User = CiDbUser
	cfg.Password = CiDbPassword
	cfg.Database = "goradd"

	db1 := pgsql.NewDB(key, "", cfg)
	opt := pgsql.DefaultOptions()
	db1.Analyze(opt)

	db.AddDatabase(db1, key)

	cfg, _ = pgx.ParseConfig("")
	cfg.Host = "localhost"
	cfg.User = CiDbUser
	cfg.Password = CiDbPassword
	cfg.Database = "goradd_unit"

	key = "goraddUnit"

	db2 := pgsql.NewDB(key, "", cfg)
	opt = pgsql.DefaultOptions()
	opt.Schemas = []string{"goradd_unit"}

	db2.Analyze(opt)
	db.AddDatabase(db2, key)

}
