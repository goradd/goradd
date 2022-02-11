package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/stringmap"
	"github.com/goradd/goradd/pkg/sys"
)

func initDatabases() {
	//if !config.Release {
	addDatabase()
	//}
}

func addDatabase() {
	var configOverrides map[string]interface{}

	// Use our own flag getter since the flags have not been read yet
	dbConfigFile, _ := sys.GetFlagString("-dbConfigFile")
	if dbConfigFile != "" {
		var err error

		configOverrides, err = stringmap.FromJsonFile(dbConfigFile)
		if err != nil {
			panic("Database configuration file error: " + err.Error())
		}
	}

	cfg := mysql.NewConfig()

	// Local development credentials
	key := "goradd"
	cfg.DBName = "goradd"
	cfg.User = "root"
	cfg.Passwd = "12345"
	cfg.ParseTime = true

	if i, ok := configOverrides[key]; ok {
		db.MysqlOverrideConfigSettings(cfg, i.(map[string]interface{}))
	}

	db1 := db.NewMysql5(key, "", cfg)

	if !config.Release {
		db1.StartProfiling()
	}

	db.AddDatabase(db1, key)
}
