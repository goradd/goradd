package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
)

// Initialize the databases that the application will use through the database query interfaces
func initDatabases() {
	// Add your database credentials below.

	// Uncomment this to make the examples work
	//addGoraddDatabase()

	// add your own development databases

}

// addGoraddDatabase adds the goradd sample database to the database list. You will need this to run some of the
// examples. You also need to install the database itself, located in the examples/db directory.
func addGoraddDatabase() {
	cfg := mysql.NewConfig()

	if !config.Release {
		// Local development credentials
		cfg.DBName = "goradd"
		cfg.User = "root"
		cfg.Passwd = "12345"
		cfg.ParseTime = true
	} else {
		// Release credentials here
	}

	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)
}
