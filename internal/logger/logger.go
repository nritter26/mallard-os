package logger

import (
	"fmt"
	"time"
)

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type Writer interface {
	Write([]byte) (int, error)
}

type Logger struct {
	started  time.Time
	writer   Writer
	minLevel Level
}

func New(writer Writer) *Logger {
	return &Logger{
		started:  time.Now(),
		writer:   writer,
		minLevel: DEBUG,
	}
}

func (l *Logger) SetLevel(level Level) {
	if level > ERROR {
		return
	}

	l.minLevel = level
}

func (l *Logger) Log(level Level, message string) {
	if l.writer == nil {
		return
	}

	if level > ERROR || level < l.minLevel {
		return
	}

	elapsed := time.Since(l.started).Milliseconds()

	var label string

	switch level {
		case DEBUG:
			label = "DEBUG"
		case INFO:
			label = "INFO"
		case WARN:
			label = "WARN"
		case ERROR:
			label = "ERROR"
	}

	_, err := fmt.Fprintf(
		l.writer,
		"[%s] %d ms %s\n",
		label,
		elapsed,
		message,
	)

	if err != nil {
		// Fallback strategy later.
	}
}

