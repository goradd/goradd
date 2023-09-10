package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var logger = loggerTest{}

type loggerTest struct {
	err string
}

func (l *loggerTest) Log(err string) {
	l.err = err
}

func setLoggers() {
	CreateDefaultLoggers()
	SetLogger(FrameworkDebugLog, &logger)
	SetLogger(InfoLog, &logger)
	SetLogger(WarningLog, &logger)
	SetLogger(ErrorLog, &logger)
	SetLogger(DebugLog, &logger)
	SetLogger(SqlLog, &logger)

	// test deleting log
	SetLogger(FrameworkInfoLog, nil)
}

func TestSetLogger(t *testing.T) {
	SetLoggingLevel(SqlLog)
	setLoggers()

	// should be empty
	FrameworkDebug("debug")
	assert.EqualValues(t, "", logger.err)

	// should be empty
	FrameworkDebugf("debugf")
	assert.EqualValues(t, "", logger.err)

	Info("debug2")
	assert.EqualValues(t, "debug2", logger.err)

	Infof("debug2f")
	assert.EqualValues(t, "debug2f", logger.err)

	Warning("debug3")
	assert.EqualValues(t, "debug3", logger.err)

	Warningf("debug3f")
	assert.EqualValues(t, "debug3f", logger.err)

	Error("debug4")
	assert.EqualValues(t, "debug4", logger.err)

	Errorf("debug4f")
	assert.EqualValues(t, "debug4f", logger.err)

	Debug("debug5")
	assert.EqualValues(t, "debug5", logger.err)

	Debugf("debug5f")
	assert.EqualValues(t, "debug5f", logger.err)

	Print(WarningLog, "warn")
	assert.EqualValues(t, "warn", logger.err)

	Printf(WarningLog, "warn2")
	assert.EqualValues(t, "warn2", logger.err)

	Sql("sql")
	assert.EqualValues(t, "sql", logger.err)

	logger.err = ""

	// should be empty
	FrameworkInfo("info")
	assert.EqualValues(t, "", logger.err)

	// should be empty
	FrameworkInfof("info")
	assert.EqualValues(t, "", logger.err)

	assert.True(t, HasLogger(SqlLog))
	assert.False(t, HasLogger(FrameworkInfoLog))
}
