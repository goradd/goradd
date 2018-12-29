package dbconfig

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
)

// Initialize the databases that the application will use through the database query interfaces
func InitDatabases() {
	cfg := mysql.NewConfig()

	if !config.Release {
		cfg.DBName = "goradd"
		cfg.User = "root"
		cfg.Passwd = "12345"
		cfg.ParseTime = true
	} else {
	}


	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)

	// add more databases


	// do this last
	db.AnalyzeDatabases()
}
