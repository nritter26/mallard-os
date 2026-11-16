package logger

import "time"

type Level uint8

const (
	INFO Level = iota
	WARN
	ERROR
	DEBUG
)

type Logger struct {
	strarted time.Time
}

func New() *Logger {
	return &Logger{
		strarted: time.Now(),
	}
}

func (l *Logger) Log(level Level, message string) {
	elapsed := time.Since(l.strarted).Milliseconds()

	var label string

	switch level {
	case INFO:
		label = "INFO"
	case WARN:
		label = "WARN"
	case ERROR:
		label = "DEBUG"
	default:
		label = "UNKNOWN"
	}

	println("[", label, "]", elapsed, "ms", message)
}


