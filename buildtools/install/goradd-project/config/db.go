package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/spekary/goradd/orm/db"
)

// Initialize the databases that the application will use through the database query interfaces
func InitDatabases() {
	cfg := mysql.NewConfig()

	// change these parameters as needed
	if !Release {
		// local development settings
		cfg.DBName = "goradd"
		cfg.User = "root"
		cfg.Passwd = "12345"
	} else {
		// settings for the release server
		cfg.DBName = "goradd"
		cfg.User = "dbuser"
		cfg.Passwd = "dbpassword"
	}


	key := "impulse"

	db1 := db.NewMysql5(key, "", cfg)

	db.AddDatabase(db1, key)

	db.AnalyzeDatabases()
}
