package config

import (
	"github.com/goradd/goradd/pkg/config"
	grlog "github.com/goradd/goradd/pkg/log"
	"log"
	"os"
)

// initLogs initializes the loggers for the various levels of logging in the app.
// You can specify a different logging behavior from here depending on whether it
// is a development or release build.
func initLogs() {
	if !config.Release {
		// Development build
		grlog.CreateDefaultLoggers()
		/* Uncomment to watch the SQL being sent to the database
		grlog.SetLogger(grlog.SqlLog, grlog.StandardLogger{Logger: log.New(os.Stderr,
			"Sql:     ", log.Ldate|log.Lmicroseconds|log.Llongfile)})
		*/
	} else {
		grlog.SetLogger(grlog.ErrorLog, grlog.StandardLogger{Logger: log.New(os.Stderr,
			"Error:     ", log.Ldate|log.Lmicroseconds|log.Llongfile)})
	}
}
