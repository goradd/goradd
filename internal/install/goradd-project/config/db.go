package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	mysql2 "github.com/goradd/goradd/pkg/orm/db/sql/mysql"
	"github.com/goradd/goradd/pkg/orm/db/sql/pgsql"
	"github.com/jackc/pgx/v5"
)

// Initialize the databases that the application will use through the database query interfaces
func initDatabases() {
	// Uncomment one of these to make the examples work
	//addMysqlGoraddDatabase()
	//addPostgresGoraddDatabase()

	// add your own development databases
	// addMyDatabase()
}

// addMysqlGoraddDatabase adds the MySQL goradd sample database to the database list.
// You will need this to run some of the examples.
// You also need to install the database itself, located in the examples/db directory.
func addMysqlGoraddDatabase() {
	cfg := mysql.NewConfig()

	// Local development credentials
	key := "goradd"
	cfg.DBName = "goradd"
	cfg.User = "root"
	cfg.Passwd = "12345"
	cfg.ParseTime = true

	// Open a database
	db1 := mysql2.NewDB(key, "", cfg)
	// Tell the database to analyze its own structure
	db1.Analyze(mysql2.DefaultOptions())

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)
}

// addPostgresGoraddDatabase adds the Postgres goradd sample database to the database list.
// You will need this to run some of the examples.
// You also need to install the database itself, located in the examples/db directory.
func addPostgresGoraddDatabase() {
	cfg, _ := pgx.ParseConfig("")

	// Local development credentials
	key := "goradd"
	cfg.Host = "localhost"
	cfg.User = "root"
	cfg.Password = "12345"
	cfg.Database = "goradd"

	db1 := pgsql.NewDB(key, "", cfg)
	opt := pgsql.DefaultOptions()
	db1.Analyze(opt)

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)
}

// addMyDatabase is a sample of how to add your own database to the database list.
// It uses a db.cfg file to hold the credentials for the deployed version of the
// app. Modify as needed.
/*
func addMyDatabase() {
	var configOverrides map[string]interface{}

	// Use our own flag getter since the flags have not been read yet
	dbConfigFile, _ := sys.GetFlagString("-dbConfigFile")

	if dbConfigFile != "" {
		var err error

		configOverrides, err = stringmap.FromJsonFile(dbConfigFile)
		if err != nil {
			panic ("Database configuration file error: " + err.Error())
		}
	}

	cfg := mysql.NewConfig()

	// Local development credentials
	key := "mydb"
	cfg.DBName = "mydbname"
	cfg.User = "root"
	cfg.Passwd = "12345"
	cfg.ParseTime = true

	// These are overridden by the database config file if one is specified on the command line.
	// The recommended way to secure database
	// credentials in your release version of the app is to store them in this json file
	// that only the launcher of the app can read. On unix, you would `chmod 400 file` to
	// set the permission of the file. The app then will read the file above and use it
	// to set the credentials. If using docker, mount the file into the docker container.
	// See goradd-project/build/docker/ for examples.
	//
	// The format of the file should be:
	// {
	// 	  "key (same as key above)": {
	//		"dbname": "database name here",
	//		"user": "database user name here",
	//		etc.
	//	  }
	// }
	// 	 See OverrideConfigSettings for more info.
	//
	if i,ok := configOverrides[key]; ok {
		mysql2.OverrideConfigSettings(cfg, i.(map[string]interface{}))
	}

	db1 := mysql2.NewDB(key, "", cfg)

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)
}*/
