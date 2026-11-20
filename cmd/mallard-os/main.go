package main

import (
	"time"

	"machine"

	"mallard-os/internal/boot"
	"mallard-os/internal/console"
	"mallard-os/internal/hal"
	"mallard-os/internal/logger"
	"mallard-os/internal/shell"
)

func main() {
	time.Sleep(3 * time.Second)

	log := logger.New(machine.Serial)

	// Initialize hardware devices
	usbConsole := hal.NewUSBConsole()

	// Create hardware-independent terminal
	terminal := console.New(usbConsole)

	// Initialize terminal logger
	terminalLog := console.NewLogger(terminal)

	// Configure status LED
	led := hal.NewGPIO(
		"status-led",
		hal.NewMachinePin(machine.GPIO21),
			   hal.PinOutput,
	)

	// Register devices with HAL
	hardware := hal.NewManager(
		usbConsole,
		led,
	)

	// Initialize boot manager
	manager := boot.NewManager(hardware)
	manager.SetLogger(log)

	// Execute boot sequence
	if err := manager.Run(); err != nil {
		log.Log(logger.ERROR, "Mallard-OS boot failed: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	// Configure platform clock
	log.Log(logger.INFO, "Configuring platform clock")

	if err := machine.SetCPUFrequency(160_000_000); err != nil {
		log.Log(logger.ERROR, "Failed to set CPU frequency")

		for {
			time.Sleep(time.Second)
		}
	}

	freq, err := machine.GetCPUFrequency()
	if err != nil {
		log.Log(logger.ERROR, "Failed to read CPU frequency")

		for {
			time.Sleep(time.Second)
		}
	}

	println("[INFO] CPU frequency:", freq)

	log.Log(logger.INFO, "Mallard-OS hardware initialized")

	// Terminal startup
	_ = terminalLog.Info("Mallard-OS starting...")
	_ = terminalLog.Success("USB initialized")
	_ = terminalLog.Warn("Running without persistent storage")

	// Initialize shell command registry
	registry := shell.NewRegistry()

	// Register built-in help command
	if err := registry.Register(shell.Command{
		Name:        "help",
		Description: "Show available commands",
		Handler: func(args []string) error {
			if err := terminal.Println("Available commands:"); err != nil {
				return err
			}

			return terminal.Println("  help    Show available commands")
		},
	}); err != nil {
		log.Log(logger.ERROR, "Failed to register help command: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	// Initialize MallardOS shell
	mallardShell := shell.NewShell(terminal, registry)

	if err := mallardShell.Start(); err != nil {
		log.Log(logger.ERROR, "Shell initialization failed: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	// Main system loop
	lastBlink := time.Now()
	ledOn := false

	for {
		// Process shell input
		if err := mallardShell.Poll(); err != nil {
			log.Log(logger.ERROR, "Shell input error: "+err.Error())
		}

		// LED heartbeat
		if time.Since(lastBlink) >= time.Second {
			ledOn = !ledOn

			if err := led.Write(ledOn); err != nil {
				log.Log(logger.ERROR, "LED toggle failed")
			}

			lastBlink = time.Now()
		}

		time.Sleep(10 * time.Millisecond)
	}
}
