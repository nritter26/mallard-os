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
	usbConsole := hal.NewUSBConsole()

	// Depend on the generic Console interface
	var console hal.Console = usbConsole

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
		log.Log(logger.ERROR, "Mallard-OS boot failed")

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

	// Console startup message
	_, _ = console.Write([]byte("\r\nMallard-OS USB Console Ready\r\n"))
	_, _ = console.Write([]byte("It's an OS with a duck...\r\n"))
	_, _ = console.Write([]byte("Waiting for USB input...\r\n"))

	// Main system loop
	lastBlink := time.Now()
	ledOn := false

	for {
		// USB receive diagnostics
		count := console.Buffered()

		if count > 0 {
			println("[DEBUG] RX bytes:", count)

			char, err := console.ReadByte()

			if err != nil {
				println("[DEBUG] Read error")
			} else {
				println("[DEBUG] Received:", char)

				_, err = console.Write([]byte{char})
				if err != nil {
					log.Log(logger.ERROR, "Console write failed")
				}
			}
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
