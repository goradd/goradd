package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/spekary/goradd/pkg/orm/db"
)

// Initialize the databases that the application will use through the database query interfaces
func InitDatabases() {
	cfg := mysql.NewConfig()

	if !Release {
		cfg.DBName = "impulse"
		cfg.User = "root"
		cfg.Passwd = "12345"
		cfg.ParseTime = true
	} else {
		cfg.DBName = "impulse"
		cfg.User = "impulse"
		cfg.Passwd = "xKMMahB35SaDnLVZ"
		cfg.ParseTime = true
	}


	key := "impulse"

	db1 := db.NewMysql5(key, "", cfg)

	if !Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)

	db.AnalyzeDatabases()
}
