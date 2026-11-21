package shell

import (
	"errors"

	"mallard-os/internal/console"
)

var (
	ErrShellNotInitialized = errors.New("shell not initialized")
	ErrShellNotStarted     = errors.New("shell not started")
)

type Shell struct {
	terminal *console.Console
	editor   *console.LineEditor
	registry *Registry
	started  bool
}

func NewShell(terminal *console.Console, registry *Registry) *Shell {
	return &Shell{
		terminal: terminal,
		editor:   console.NewLineEditor(terminal),
		registry: registry,
	}
}

// Start activates the shell and displays its welcome message
func (s *Shell) Start() error {
	if s.terminal == nil || s.editor == nil || s.registry == nil {
		return ErrShellNotInitialized
	}

	if s.started {
		return nil
	}

	banner := []string{
		"       __",
		"  ___( o)>",
		"  \\ <_. )",
		"   `---'",
		"",
		"  MALLARD-OS v0.1",
		"  Embedded Operating Environment",
		"",
	}

	for _, line := range banner {
		if err := s.terminal.Println(line); err != nil {
			return err
		}
	}

	if err := s.terminal.Println("Type 'help' to see available commands."); err != nil {
		return err
	}

	if err := s.terminal.Println(""); err != nil {
		return err
	}

	s.started = true

	return s.printPrompt()
}

// Poll processes available terminal input without blocking
func (s *Shell) Poll() error {
	if !s.started {
		return ErrShellNotStarted
	}

	line, ready, err := s.editor.Poll()
	if err != nil {
		return err
	}

	if !ready {
		return nil
	}

	command, args, err := Parse(line)
	if err != nil {
		if printErr := s.terminal.Println("parse error: " + err.Error()); printErr != nil {
			return printErr
		}

		return s.printPrompt()
	}

	// Ignore empty input
	if command == "" {
		return s.printPrompt()
	}

	err = s.registry.Execute(command, args, s.terminal)
	if err != nil {
		if errors.Is(err, ErrCommandNotFound) {
			err = s.terminal.Println("unknown command: " + command)
		} else {
			err = s.terminal.Println("command error: " + err.Error())
		}

		if err != nil {
			return err
		}
	}

	return s.printPrompt()
}

func (s *Shell) printPrompt() error {
	return s.terminal.Print("mallard> ")
}
