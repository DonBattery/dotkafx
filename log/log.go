package log

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	DebugLevel    LogLevel = 0
	InfoLevel     LogLevel = 1
	WarnLevel     LogLevel = 2
	ErrorLevel    LogLevel = 3
	FatalLevel    LogLevel = 4
	ShutdownLevel LogLevel = 5
)

var LoggingLevel = InfoLevel

func levelToString(level LogLevel) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info "
	case WarnLevel:
		return "warn "
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case ShutdownLevel:
		return "shutdown"
	default:
		return "unknown"
	}
}

func Entry(level LogLevel, format string, args ...any) {
	if level < LoggingLevel {
		return
	}

	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}

	fmt.Printf("%s %s - %s\n", time.Now().Format("2006-01-02 15:04:05"), levelToString(level), format)

	if level == FatalLevel {
		os.Exit(1)
	}

	if level == ShutdownLevel {
		os.Exit(0)
	}
}

func Debug(format string, args ...any) {
	Entry(DebugLevel, format, args...)
}

func Info(format string, args ...any) {
	Entry(InfoLevel, format, args...)
}

func Warn(format string, args ...any) {
	Entry(WarnLevel, format, args...)
}

func Error(format string, args ...any) {
	Entry(ErrorLevel, format, args...)
}

func Fatal(format string, args ...any) {
	Entry(FatalLevel, format, args...)
}

func Shutdown(format string, args ...any) {
	Entry(ShutdownLevel, format, args...)
}
