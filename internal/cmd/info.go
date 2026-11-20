package cmd

import "mallard-os/internal/shell"

func Info(io shell.CommandIO, args []string) error {
	if err := io.Println("Mallard-OS v0.1"); err != nil {
		return err
	}

	if err := io.Println("Embedded Operating Environment"); err != nil {
		return err
	}

	if err := io.Println("Architecture: ESP32-S3"); err != nil {
		return err
	}

	return io.Println("Runtime: TinyGo")
}
