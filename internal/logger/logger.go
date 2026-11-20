package logger

import "time"

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
	started time.Time
	writer Writer
}

func New(writer Writer) *Logger {
	return &Logger{
		started: time.Now(),
		writer: writer,
	}
}

func (l *Logger) Log(level Level, message string) {
	elapsed := time.Since(l.started).Milliseconds()

	var label string

	switch level {
	case INFO:
		label = "INFO"
	case WARN:
		label = "WARN"
	case ERROR:
		label = "ERROR"
	case DEBUG:
		label = "DEBUG"
	default:
		label = "UNKNOWN"
	}

	_, err := fmt.Fprintf(l.writer, "[%s] %d ms %s\n", label, elapsed, message)
	if err != nil {

}


