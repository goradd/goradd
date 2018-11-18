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
		cfg.ParseTime = true
		// Set the Loc to the timezone setting of your MySQL server. The default is UTC, so only do the below if you server is set to a different time zone. UTC is best for portability of the data.
		//cfg.Loc,_ = time.LoadLocation("America/Los_Angeles")
	} else {
		// settings for the release server
		cfg.DBName = "goradd"
		cfg.User = "dbuser"
		cfg.Passwd = "dbpassword"
		cfg.ParseTime = true
		// Set the Loc to the timezone setting of your MySQL server. The default is UTC, so only do the below if you server is set to a different time zone. UTC is best for portability of the data.
		//cfg.Loc,_ = time.LoadLocation("America/Los_Angeles")
	}


	key := "impulse"

	db1 := db.NewMysql5(key, "", cfg)

	//db1.StartProfiling() turn this on to profile your sql queries

	db.AddDatabase(db1, key)

	db.AnalyzeDatabases()
}
