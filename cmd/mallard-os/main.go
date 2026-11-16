package main

import (
	"time"

	"mallard-os/internal/boot"
	"mallard-os/internal/hal"
	"mallard-os/internal/logger"

	"machine"
)

func main() {
	time.Sleep(3 * time.Second)

	log := logger.New()

	// Initialize hardware devices
	console := hal.NewUSBConsole()

	led := hal.NewGPIO(
		"status-led",
		hal.NewMachinePin(machine.GPIO21),
			   hal.PinOutput,
	)

	// Register devices with HAL
	hardware := hal.NewManager(
		console,
		led,
	)

	// Initialize boot manager
	manager := boot.NewManager(hardware)
	manager.SetLogger(log)

	// Execute boot sequence
	if err := manager.Run(); err != nil {
		log.Log(logger.ERROR, "MallardOS boot failed")

		for {
		}
	}

	// Configure platform clock
	log.Log(logger.INFO, "Configuring platform clock")

	if err := machine.SetCPUFrequency(160_000_000); err != nil {
		log.Log(logger.ERROR, "Failed to set CPU frequency")

		for {
		}
	}

	freq, err := machine.GetCPUFrequency()
	if err != nil {
		log.Log(logger.ERROR, "Failed to read CPU frequency")

		for {
		}
	}

	println("[INFO] CPU frequency:", freq)

	log.Log(logger.INFO, "MallardOS hardware initialized")

	// Main system loop
	for {
		if err := led.Write(true); err != nil {
			log.Log(logger.ERROR, "LED ON failed")
		}

		time.Sleep(time.Second)

		if err := led.Write(false); err != nil {
			log.Log(logger.ERROR, "LED OFF failed")
		}

		time.Sleep(time.Second)

		println("[INFO] Mallard-OS alive")
	}
}
