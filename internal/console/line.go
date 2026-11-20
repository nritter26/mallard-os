package console

import "errors"

const MaxLineLength = 128

var (
	ErrLineEditorNotInitialized = errors.New("line editor not initialized")
)

// LineEditor manages interactive terminal input
type LineEditor struct {
	terminal *Console

	buffer [MaxLineLength]byte
	length int

	skipLF bool
}

// NewLineEditor creates a new line editor
func NewLineEditor(terminal *Console) *LineEditor {
	return &LineEditor{
		terminal: terminal,
	}
}

// Poll processes available input without blocking
//
// Returns:
//   - line: completed input
//   - ready: whether a complete line is available
//   - error: input/output error
func (e *LineEditor) Poll() (string, bool, error) {
	if e.terminal == nil {
		return "", false, ErrLineEditorNotInitialized
	}

	for e.terminal.Buffered() > 0 {
		char, err := e.terminal.ReadByte()
		if err != nil {
			return "", false, err
		}

		// Handle CRLF sequences
		if e.skipLF {
			e.skipLF = false

			if char == '\n' {
				continue
			}
		}

		switch char {

			case '\r':
				e.skipLF = true
				return e.finishLine()

			case '\n':
				return e.finishLine()

				// Backspace and Delete
			case '\b', 127:
				if e.length > 0 {
					e.length--

					if err := e.terminal.Print("\b \b"); err != nil {
						return "", false, err
					}
				}

			default:
				// Accept printable ASCII characters
				if char >= 32 && char <= 126 {
					if e.length >= MaxLineLength {
						if err := e.terminal.Print("\a"); err != nil {
							return "", false, err
						}
						continue
					}

					e.buffer[e.length] = char
					e.length++

					if _, err := e.terminal.Write([]byte{char}); err != nil {
						return "", false, err
					}
				}
		}
	}

	return "", false, nil
}

// finishLine returns the completed input line
func (e *LineEditor) finishLine() (string, bool, error) {
	if err := e.terminal.Print("\r\n"); err != nil {
		return "", false, err
	}

	line := string(e.buffer[:e.length])

	e.length = 0

	return line, true, nil
}
