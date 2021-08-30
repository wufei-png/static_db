package log

type LogLevel int

const (
	TRACE LogLevel = 0
	DEBUG          = 1
	INFO           = 2
	WARN           = 3
	ERROR          = 4
	FATAL          = 5
)

type Logger interface {
	Logf(level LogLevel, fmt string, args ...interface{})
	Log(level LogLevel, args ...interface{})
}

var DefaultStandardLogger = &StandardLogger{Level: INFO}
var DefaultLogger Logger = DefaultStandardLogger

func Logf(level LogLevel, fmt string, args ...interface{}) {
	DefaultLogger.Logf(level, fmt, args...)
}

func Log(level LogLevel, args ...interface{}) {
	DefaultLogger.Log(level, args...)
}
