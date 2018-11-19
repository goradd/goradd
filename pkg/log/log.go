// Package log provides logging support for the application and framework developers
package log

import (
	"fmt"
	"log"
	"goradd-project/config"
)

const FrameworkDebugLog = 1 // Should only exist when debugging style logs are required for debugging the goradd framework itself.
const InfoLog = 2           // Info log is designed to always exist. These would be messages we only need to check periodically to know the system is working correctly.
const WarningLog = 3        // Should be sent to sysop periodically (daily perhaps?). Would generally be issues involving low resources.
const ErrorLog = 4          // Should be sent to sysop immediately
const DebugLog = 10         // Debug log for the developer's application, separate from the goradd framework debug log

var Loggers = map[int]*log.Logger{}

func HasLogger(id int) bool {
	_, ok := Loggers[id]
	return ok
}

func Info(v ...interface{}) {
	_print(InfoLog, v...)
}

func Infof(format string, v ...interface{}) {
	_printf(InfoLog, format, v...)
}


func FrameworkDebug(v ...interface{}) {
	if config.Debug {
		_print(FrameworkDebugLog, v...)
	}
}

func FrameworkDebugf(format string, v ...interface{}) {
	if config.Debug {
		_printf(FrameworkDebugLog, format, v...)
	}
}

func Warning(v ...interface{}) {
	_print(WarningLog, v...)
}

func Warningf(format string, v ...interface{}) {
	_printf(WarningLog, format, v...)
}

func Error(v ...interface{}) {
	_print(ErrorLog, v...)
}

func Errorf(format string, v ...interface{}) {
	_printf(ErrorLog, format, v...)
}

// Debug is for application debugging logging
func Debug(v ...interface{}) {
	if config.Debug {
		_print(DebugLog, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if config.Debug {
		_printf(DebugLog, format, v...)
	}
}

func Print(logType int, v ...interface{}) {
	_print(logType, v)
}

func Printf(logType int, format string, v ...interface{}) {
	_printf(logType, format, v)
}


func _print(logType int, v ...interface{}) {
	if l, ok := Loggers[logType]; ok {
		l.Output(3, fmt.Sprint(v...))
	}
}

func _printf(logType int, format string, v ...interface{}) {
	if l, ok := Loggers[logType]; ok {
		l.Output(3, fmt.Sprintf(format, v...))
	}
}
