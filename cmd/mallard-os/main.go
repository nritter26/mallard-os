package main

import (
	"time"
	"machine"

	commands "mallard-os/internal/commands"
	mallardruntime "mallard-os/internal/runtime"
	"mallard-os/internal/kernel"
	"mallard-os/internal/console"
	"mallard-os/internal/hal"
	"mallard-os/internal/logger"
	"mallard-os/internal/shell"
)

func main() {

	// Initialize run time clock
	mallardruntime.Initialize()

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

	// Initialize Mallard-OS kernel
	mallardKernel := kernel.New()

	// Register hardware subsystem
	if err := mallardKernel.Register(hardware); err != nil {
		log.Log(logger.ERROR, "Failed to register HAL: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	// Initialize kernel and registered subsystems
	log.Log(logger.INFO, "Initializing Mallard-OS kernel")

	if err := mallardKernel.Initialize(); err != nil {
		log.Log(logger.ERROR, "Kernel initialization failed: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	log.Log(logger.INFO, "Mallard-OS kernel initialized")

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

	// Register built-in commands
	if err := commands.Register(registry); err != nil {
		log.Log(logger.ERROR, "Failed to register built-in commands: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	// Initialize Mallard-OS shell
	mallardShell := shell.NewShell(terminal, registry)

	if err := mallardShell.Start(); err != nil {
		log.Log(logger.ERROR, "Shell initialization failed: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	// Start kernel
	if err := mallardKernel.Start(); err != nil {
		log.Log(logger.ERROR, "Kernel startup failed: "+err.Error())

		for {
			time.Sleep(time.Second)
		}
	}

	log.Log(logger.INFO, "Mallard-OS kernel running")

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
