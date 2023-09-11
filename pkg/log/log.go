// Package log provides logging support for your application and the goradd framework developers.
//
// The logging package is fashioned after the slog package that was released in Go 1.21, but
// has some notable differences. In particular, each log level is not just a level, but a separate
// logger, so that you can create very different application responses depending on the log level
// of the event being recorded. For example, if you would like to email certain logging events to
// the sysop, you can create a particular logger for that log level.
//
// While developing the slog package, the Go developers did some performance testing, and
// found that memory allocation was the slowest part of logging. So, they created *Attrs calls
// that lazy-loaded the attributes and strings, essentially waiting to make sure that a Log call
// was going to be used before turning an attribute into a string. This principal is attempted to
// be used here as well, but will require some transition time.
//
// The default main.go file has a command line flag that calls SetLoggingLevel at startup time.
package log

import (
	"fmt"
	"log"
	"os"

	"github.com/goradd/goradd/pkg/config"
)

// The log constants represent both a separate logger, and a log level.
// Set the log level to turn on or off specific loggers.
// These logging levels correspond roughly to slog's logging levels.
const (
	FrameworkDebugLog = -12 // Used by framework developers for debugging the framework
	SqlLog            = -8  // Outputs the SQL from all SQL queries. This should be used carefully since sensitive information may appear in the logs.
	DebugLog          = -4  // For use by the application for application level debugging
	FrameworkInfoLog  = 0   // Used by the framework for logging normal application progress, typically a log of http calls. This is the default logging level.
	InfoLog           = 1   // For use by the application to do any other info level logging. Set this logging level to turn off the frameworks info logs so that only the application info logs are reported.
	WarningLog        = 4   // For use by the framework and the application for information that would need to be sent to the sysop periodically. Reports possible performance issues.
	ErrorLog          = 8   // For use by the framework and application for information that should be sent to the sysop immediately.
)

// loggingLevel is the current logging level. It should only be changed at system startup.
var loggingLevel = FrameworkInfoLog

// loggers is the global map of loggers in use.
var loggers = make(map[int]LoggerI)

// SetLoggingLevel sets the current logger level. The current implementation is only suitable
// to be set at startup time and not while the application is running.
func SetLoggingLevel(l int) {
	loggingLevel = l
}

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
		panic("Logging error: " + err.Error())
	}
}

func (l EmailLogger) Log(out string) {
	if err := l.Output(4, out); err != nil {
		panic("Logging error: " + err.Error())
	}

	// TODO: Create emailer
	// email.ErrorSend(l.EmailAddresses, "Error", out)
}

// HasLogger returns true if the given logger is available.
func HasLogger(id int) bool {
	_, ok := loggers[id]
	return ok
}

// SetLogger sets the given logger id to the given logger.
//
// Use this function to set up the loggers at startup time.
// Pass nil for a logger to delete it.
func SetLogger(id int, l LoggerI) {
	if l == nil {
		delete(loggers, id)
	} else {
		loggers[id] = l
	}
}

// Info prints information to the Info logger.
func Info(v ...interface{}) {
	_print(InfoLog, v...)
}

// Infof prints formatted information to the InfoLog logger.
func Infof(format string, v ...interface{}) {
	_printf(InfoLog, format, v...)
}

// FrameworkDebug is used by the framework to log debugging information.
// This is mostly for development of the framework, but it can also be used to track down problems
// in your own code.
func FrameworkDebug(v ...interface{}) {
	_print(FrameworkDebugLog, v...)
}

// FrameworkDebugf is used by the framework to log formatted debugging information.
func FrameworkDebugf(format string, v ...interface{}) {
	_printf(FrameworkDebugLog, format, v...)
}

// FrameworkInfo is used by the framework to log operational information.
func FrameworkInfo(v ...interface{}) {
	_print(FrameworkInfoLog, v...)
}

// FrameworkInfof is used by the framework to log formatted operational information.
func FrameworkInfof(format string, v ...interface{}) {
	_printf(FrameworkInfoLog, format, v...)
}

// Sql outputs the sql sent to the database driver.
func Sql(v ...interface{}) {
	_print(SqlLog, v...)
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

// Errorf prints formmated information to the Error logger.
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
	if loggingLevel <= logType {
		if l, ok := loggers[logType]; ok && l != nil {
			l.Log(fmt.Sprint(v...))
		}
	}
}

func _printf(logType int, format string, v ...interface{}) {
	if loggingLevel <= logType {
		if l, ok := loggers[logType]; ok && l != nil {
			l.Log(fmt.Sprintf(format, v...))
		}
	}
}

// CreateDefaultLoggers creates default loggers for the application.
// After calling this, you can replace the loggers with your own by calling SetLogger, and
// add additional loggers to the logging array, or remove ones you don't use.
func CreateDefaultLoggers() {
	// make these strings the same size to improve the look of the log
	loggers[FrameworkDebugLog] = StandardLogger{log.New(os.Stdout,
		"Framework: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)}
	loggers[SqlLog] = StandardLogger{log.New(os.Stdout,
		"Framework: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)}
	loggers[FrameworkInfoLog] = StandardLogger{log.New(os.Stdout,
		"Framework: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)}
	loggers[InfoLog] = StandardLogger{log.New(os.Stdout,
		"Info:      ", log.Ldate|log.Lmicroseconds)}
	loggers[DebugLog] = StandardLogger{log.New(os.Stdout,
		"Debug:     ", log.Ldate|log.Lmicroseconds|log.Lshortfile)}
	loggers[WarningLog] = StandardLogger{log.New(os.Stderr,
		"Warning:   ", log.Ldate|log.Lmicroseconds|log.Llongfile)}
	loggers[ErrorLog] = StandardLogger{log.New(os.Stderr,
		"Error:     ", log.Ldate|log.Lmicroseconds|log.Llongfile)}
}
