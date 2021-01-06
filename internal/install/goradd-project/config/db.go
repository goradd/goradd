package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
)

// Initialize the databases that the application will use through the database query interfaces
func initDatabases() {
	// Uncomment this to make the examples work
	//addGoraddDatabase()

	// add your own development databases
	// addMyDatabase()
}

// addGoraddDatabase adds the goradd sample database to the database list. You will need this to run some of the
// examples. You also need to install the database itself, located in the examples/db directory.
func addGoraddDatabase() {
	cfg := mysql.NewConfig()

	// Local development credentials
	key := "goradd"
	cfg.DBName = "goradd"
	cfg.User = "root"
	cfg.Passwd = "12345"
	cfg.ParseTime = true

	db1 := db.NewMysql5(key, "", cfg)

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
	// that only the launcher of the app can read. On unix, you would `chmod 400 pw_file` to
	// set the permission of the file. The app then will read the file above and use it
	// to set the credentials. If using docker, mount the file into the docker container.
	//
	// The format of the file should be:
	// {
	// 	  "put_db_key_here": {
	//		"dbname": "database name here",
	//		"user": "database user name here",
	//		etc.
	// 	 See MysqlOverrideConfigSettings for more info.
	//
	if i,ok := configOverrides[key]; ok {
		db.MysqlOverrideConfigSettings(cfg, i.(map[string]interface{}))
	}

	db1 := db.NewMysql5(key, "", cfg)

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)
}*/
