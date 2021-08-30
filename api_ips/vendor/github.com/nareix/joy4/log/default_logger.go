package log

import "log"

type StandardLogger struct {
	Level LogLevel
}

func levelToName(l LogLevel) string {
	switch l {
	case TRACE:
		return "[TRACE]"
	case DEBUG:
		return "[DEBUG]"
	case INFO:
		return "[INFO ]"
	case WARN:
		return "[WARN ]"
	case ERROR:
		return "[ERROR]"
	case FATAL:
		return "[FATAL]"
	}
	return "[UNKNOWN]"
}

func (l *StandardLogger) Logf(level LogLevel, fmt string, args ...interface{}) {
	if level < l.Level {
		return
	}
	if level >= FATAL {
		log.Fatalf(levelToName(level)+fmt, args...)
	} else {
		log.Printf(levelToName(level)+fmt, args...)
	}
}

func (l *StandardLogger) Log(level LogLevel, args ...interface{}) {
	if level < l.Level {
		return
	}
	newargs := make([]interface{}, 0, len(args)+1)
	newargs = append(newargs, levelToName(level))
	newargs = append(newargs, args...)
	if level >= FATAL {
		log.Fatal(newargs...)
	} else {
		log.Print(newargs...)
	}
}
