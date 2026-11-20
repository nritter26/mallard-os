package shell

import "errors"

var (
	ErrUnclosedQuote = errors.New("unclosed quote")
	ErrInvalidEscape = errors.New("invalid escape")
)

// Parse converts a raw input line into a command and its arguments
func Parse(input string) (string, []string, error) {
	var args []string
	var current []byte

	var quote rune
	escaped := false
	inToken := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		if escaped {
			current = append(current, c)
			escaped = false
			inToken = true
			continue
		}

		if c == '\\' {
			escaped = true
			inToken = true
			continue
		}

		if quote != 0 {
			if rune(c) == quote {
				quote = 0
			} else {
				current = append(current, c)
			}
			continue
		}

		switch c {
			case '"', '\'':
				quote = rune(c)
				inToken = true

			case ' ', '\t':
				if inToken {
					args = append(args, string(current))
					current = nil
					inToken = false
				}

			default:
				current = append(current, c)
				inToken = true
		}
	}

	if escaped {
		return "", nil, ErrInvalidEscape
	}

	if quote != 0 {
		return "", nil, ErrUnclosedQuote
	}

	if inToken {
		args = append(args, string(current))
	}

	if len(args) == 0 {
		return "", nil, nil
	}

	return args[0], args[1:], nil
}
