package shell

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrDuplicateCommand = errors.New("command already registered")
	ErrCommandNotFound = errors.New("command not found")
	ErrInvalidCommand = errors.New("invalid command")
)

type CommandHandler func(io CommandIO, args []string) error

type Command struct {
	Name string
	Description string
	Handler CommandHandler
}

type Registry struct {
	commands map[string]Command
}

type CommandIO interface {
	Print(text string) error
	Println(text string) error
}
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register adds a command to the registry
func (r *Registry) Register(cmd Command) error {
	if cmd.Name == "" || cmd.Handler == nil {
		return ErrInvalidCommand
	}

	if strings.ContainsAny(cmd.Name, "\t\r\n") {
		return ErrInvalidCommand
	}

	if _, exists := r.commands[cmd.Name]; exists {
		return ErrDuplicateCommand
	}

	r.commands[cmd.Name] = cmd
	return nil
}

//Lookup retrives a registered command
func (r *Registry) Lookup(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// Execute runs a registered command
func (r *Registry) Execute(name string, args []string, io CommandIO) error {
	cmd, exists := r.Lookup(name)
	if !exists {
		return ErrCommandNotFound
	}

	return cmd.Handler(io, args)
}

// Commands returns all registered commands in alphabetical order
func (r *Registry) Commands() []Command {
	commands := make([]Command, 0, len(r.commands))

	for _, command := range r.commands {
		commands = append(commands, command)
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}
