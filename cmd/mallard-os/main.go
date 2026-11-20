package main

import (
	"time"

	"machine"

	"mallard-os/internal/boot"
	"mallard-os/internal/console"
	"mallard-os/internal/hal"
	"mallard-os/internal/logger"
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

	// Initialize line editor
	editor := console.NewLineEditor(terminal)

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

	log.Log(logger.INFO, "MallardOS hardware initialized")

	// Terminal startup
	_ = terminalLog.Info("MallardOS starting...")
	_ = terminalLog.Success("USB initialized")
	_ = terminalLog.Warn("Running without persistent storage")
	_ = terminal.Println("It's an OS with a duck...")
	_ = terminal.Println("Waiting for input...")
	_ = terminal.Print("> ")

	// Main system loop
	lastBlink := time.Now()
	ledOn := false

	for {
		// Process terminal input
		line, ready, err := editor.Poll()

		if err != nil {
			log.Log(logger.ERROR, "Terminal input error")
		}

		if ready {
			_ = terminalLog.Info("Received: " + line)
			_ = terminal.Print("> ")
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
