package cmd

import "mallard-os/internal/shell"

func Register(registry *shell.Registry) error {
	commands := []shell.Command{
		{
			Name:        "help",
			Description: "Display available commands",
			Handler:     Help(registry),
		},
		{
			Name:        "info",
			Description: "Display Mallard-OS system information",
			Handler:     Info,
		},
	}

	for _, command := range commands {
		if err := registry.Register(command); err != nil {
			return err
		}
	}

	return nil
}
