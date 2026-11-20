package cmd

import "mallard-os/internal/shell"

// Help returns a handler that displays registered commands
func Help(registry *shell.Registry) shell.CommandHandler {
	return func(io shell.CommandIO, args []string) error {
		if err := io.Println("Mallard-OS Commands:"); err != nil {
			return err
		}

		for _, command := range registry.Commands() {
			line := command.Name + " - " + command.Description

			if err := io.Println(line); err != nil {
				return err
			}
		}

		return nil
	}
}
