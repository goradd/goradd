// Package log provides logging support for the application and framework developers
package log

import (
	"fmt"
	"log"
	"goradd-project/config"
	"os"
)

const FrameworkDebugLog = 1 // Should only exist when debugging style logs are required for debugging the goradd framework itself.
const InfoLog = 2           // Info log is designed to always exist. These would be messages we only need to check periodically to know the system is working correctly.
const WarningLog = 3        // Should be sent to sysop periodically (daily perhaps?). Would generally be issues involving low resources.
const ErrorLog = 4          // Should be sent to sysop immediately
const DebugLog = 10         // Debug log for the developer's application, separate from the goradd framework debug log

type LoggerI interface {
	Log(err string)
}

type StandardLogger struct {
	*log.Logger
}

type EmailLogger struct {
	StandardLogger
	EmailAddresses []string
}

func (l StandardLogger) Log(out string) {
	l.Output(4, out)
}

func (l EmailLogger) Log(out string) {
	l.Output(4, out)

	// TODO: Create emailer
	// email.ErrorSend(l.EmailAddresses, "Error", out)
}



var Loggers = map[int]LoggerI{}

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
	_print(logType, v...)
}

func Printf(logType int, format string, v ...interface{}) {
	_printf(logType, format, v...)
}


func _print(logType int, v ...interface{}) {
	if l, ok := Loggers[logType]; ok {
		l.Log(fmt.Sprint(v...))
	}
}

func _printf(logType int, format string, v ...interface{}) {
	if l, ok := Loggers[logType]; ok {
		l.Log(fmt.Sprintf(format, v...))
	}
}


// Create's default loggers for the application. After calling this, you can replace the loggers with your own, and
// add additional loggers to the logging array.
func CreateDefaultLoggers() {
	// make these strings the same size to improve the look of the log
	Loggers[FrameworkDebugLog] = StandardLogger{log.New(os.Stdout,
		"Framework: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)}
	Loggers[InfoLog] = StandardLogger{log.New(os.Stdout,
		"Info:      ", log.Ldate|log.Lmicroseconds)}
	Loggers[DebugLog] = StandardLogger{log.New(os.Stdout,
		"Debug:     ", log.Ldate|log.Lmicroseconds|log.Lshortfile)}
	Loggers[WarningLog] = StandardLogger{log.New(os.Stderr,
		"Warning:   ", log.Ldate|log.Lmicroseconds|log.Llongfile)}
	Loggers[ErrorLog] = StandardLogger{log.New(os.Stderr,
		"Error:     ", log.Ldate|log.Lmicroseconds|log.Llongfile)}
}
