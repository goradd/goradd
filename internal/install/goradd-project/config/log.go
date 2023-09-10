package config

import (
	"github.com/goradd/goradd/pkg/config"
	grlog "github.com/goradd/goradd/pkg/log"
)

// initLogs initializes the loggers for the various levels of logging in the app.
// You can specify a different logging behavior from here depending on whether it
// is a development or release build.
func initLogs() {
	if !config.Release {
		grlog.CreateDefaultLoggers()
	} else {
		grlog.CreateDefaultLoggers()
	}
}
