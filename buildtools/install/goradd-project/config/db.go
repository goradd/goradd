package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/spekary/goradd/orm/db"
)

// Initialize the databases that the application will use through the database query interfaces
func InitDatabases() {
	// This installs the test database. Replace this with your own database.
	cfg := mysql.NewConfig()
	cfg.DBName = "goradd"
	//cfg.DBName = "test"
	cfg.User = "root"
	cfg.Passwd = "12345"
	key := "goradd"

	db1 := db.NewMysql5(key, "", cfg)
	//db1.StartProfiling() // comment out to stop profiling
	db.AddDatabase(db1, key)

	/* If using 2 databases, this is how you would do a 2nd database.
	cfg = mysql.NewConfig()
	cfg.DBName = "db2"
	//cfg.DBName = "test"
	cfg.User = "root"
	cfg.Passwd = "12345"
	key = "db2"

	db2 := db.NewMysql5(key, "", cfg)
	db2.StartProfiling() // comment out to stop profiling
	db.AddDatabase(db21, key)
*/

	db.AnalyzeDatabases()
}
