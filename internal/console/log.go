package console

// Logger formats messages for human-readable terminal output
type Logger struct {
	terminal *Console
}

// NewLogger creates a terminal logger
func NewLogger(terminal *Console) *Logger {
	return &Logger{
		terminal: terminal,
	}
}

// Info displays informational messages
func (l *Logger) Info(message string) error {
	return l.terminal.Println("[INFO] " + message)
}

// Warn displays warning messages
func (l *Logger) Warn(message string) error {
	return l.terminal.Println("[WARN] " + message)
}

// Error displays error messages
func (l *Logger) Error(message string) error {
	return l.terminal.Println("[ERROR] " + message)
}

// Success displays successful operations
func (l *Logger) Success(message string) error {
	return l.terminal.Println("[ OK ] " + message)
}
