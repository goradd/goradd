// Package log provides logging support for the application and framework developers
package log

import (
	"fmt"
	"github.com/goradd/goradd/pkg/config"
	"log"
	"os"
)

const FrameworkDebugLog = 1 // Should only exist when debugging style logs are required for debugging the goradd framework itself.
const InfoLog = 2           // Info log is designed to always exist. These would be messages we only need to check periodically to know the system is working correctly.
const WarningLog = 3        // Should be sent to sysop periodically (daily perhaps?). Would generally be issues involving low resources.
const ErrorLog = 4          // Should be sent to sysop immediately
const DebugLog = 10         // Debug log for the developer's application, separate from the goradd framework debug log

// LoggerI is the interface for all loggers.
type LoggerI interface {
	Log(err string)
}

// StandardLogger reuses the GO default logger
type StandardLogger struct {
	*log.Logger
}

// EmailLogger is a logger that will email logging activity to the emails provided.
type EmailLogger struct {
	StandardLogger
	EmailAddresses []string
}

func (l StandardLogger) Log(out string) {
	if err := l.Output(4, out); err != nil {
		panic ("Logging error: " + err.Error())
	}
}

func (l EmailLogger) Log(out string) {
	if err := l.Output(4, out); err != nil {
		panic ("Logging error: " + err.Error())
	}

	// TODO: Create emailer
	// email.ErrorSend(l.EmailAddresses, "Error", out)
}

// Loggers is the global map of loggers in use.
var Loggers = map[int]LoggerI{}

// HasLogger returns true if the given logger is available.
func HasLogger(id int) bool {
	_, ok := Loggers[id]
	return ok
}

func SetLogger(id int, l LoggerI) {
	if l == nil {
		delete (Loggers, id)
	} else {
		Loggers[id] = l
	}
}

// Info prints information to the Info logger.
func Info(v ...interface{}) {
	_print(InfoLog, v...)
}

// Infof prints formatted information to the Info logger.
func Infof(format string, v ...interface{}) {
	_printf(InfoLog, format, v...)
}

// FrameworkDebug is used by the framework to log debugging information.
// This is mostly for development of the framework, but it can also be used to track down problems
// in your own code. It becomes a no-op in the release build.
func FrameworkDebug(v ...interface{}) {
	if config.Debug {
		_print(FrameworkDebugLog, v...)
	}
}

// FrameworkDebugf is used by the framework to log formatted debugging information.
// It becomes a no-op in the release build.
func FrameworkDebugf(format string, v ...interface{}) {
	if config.Debug {
		_printf(FrameworkDebugLog, format, v...)
	}
}

// Warning prints to the Warning logger.
func Warning(v ...interface{}) {
	_print(WarningLog, v...)
}

// Warningf prints formatted information to the Warning logger.
func Warningf(format string, v ...interface{}) {
	_printf(WarningLog, format, v...)
}

// Error prints to the Error logger.
func Error(v ...interface{}) {
	_print(ErrorLog, v...)
}

// Error prints formmated information to the Error logger.
func Errorf(format string, v ...interface{}) {
	_printf(ErrorLog, format, v...)
}

// Debug is for application debugging logging. It becomes a no-op in the release build.
func Debug(v ...interface{}) {
	if config.Debug {
		_print(DebugLog, v...)
	}
}

// Debugf prints formatted information to the application Debug logger.
// It becomes a no-op in the release build.
func Debugf(format string, v ...interface{}) {
	if config.Debug {
		_printf(DebugLog, format, v...)
	}
}

// Print prints information to the given logger.
// You can use this to set up your own custom logger.
func Print(logType int, v ...interface{}) {
	_print(logType, v...)
}

// Printf prints formatted information to the given logger.
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


// CreateDefaultLoggers create's default loggers for the application.
// After calling this, you can replace the loggers with your own, and
// add additional loggers to the logging array, or remove ones you don't use.
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
