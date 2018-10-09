// Package log provides logging support for the application and
package log

import (
	"log"
	"goradd-project/config"
)

// Logging support

const FrameworkDebugLog = 1 // Should only exist when debugging style logs are required for debugging the goradd framework itself.
const InfoLog = 2           // Info log is designed to always exist. These would be messages we only need to check periodically to know the system is working correctly.
const WarningLog = 3        // Should be sent to sysop periodically (daily perhaps?). Would generally be issues involving low resources.
const ErrorLog = 4          // Should be sent to sysop immediately
const DebugLog = 10         // Debug log for the developer's application, separate from the goradd debug log

var Loggers = map[int]*log.Logger{}

func HasLogger(id int) bool {
	_, ok := Loggers[id]
	return ok
}

func Info(v ...interface{}) {
	Print(InfoLog, v...)
}

func Infof(format string, v ...interface{}) {
	Printf(InfoLog, format, v...)
}


func FrameworkDebug(v ...interface{}) {
	if config.Debug {
		Print(FrameworkDebugLog, v...)
	}
}

func FrameworkDebugf(format string, v ...interface{}) {
	if config.Debug {
		Printf(FrameworkDebugLog, format, v...)
	}
}

func Warning(v ...interface{}) {
	Print(WarningLog, v...)
}

func Warningf(format string, v ...interface{}) {
	Printf(WarningLog, format, v...)
}

func Error(v ...interface{}) {
	Print(ErrorLog, v...)
}

func Errorf(format string, v ...interface{}) {
	Printf(ErrorLog, format, v...)
}

// Debug is for application debugging logging
func Debug(v ...interface{}) {
	if config.Debug {
		Print(DebugLog, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if config.Debug {
		Printf(DebugLog, format, v...)
	}
}

func Print(logType int, v ...interface{}) {
	if l, ok := Loggers[logType]; ok {
		l.Print(v...)
	}
}

func Printf(logType int, format string, v ...interface{}) {
	if l, ok := Loggers[logType]; ok {
		l.Printf(format, v...)
	}
}
