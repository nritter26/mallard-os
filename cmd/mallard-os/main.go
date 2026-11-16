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

	console := hal.NewUSBConsole()

	hardware := hal.NewManager(
		console,
	)

	manager := boot.NewManager(hardware)
	manager.SetLogger(log)

	if err := manager.Run(); err != nil {
		log.Log(logger.ERROR, "MallardOS boot failed")
		for {
		}
	}

	log.Log(logger.INFO, "Starting GPIO blink test")

	if err := machine.SetCPUFrequency(160_000_000); err != nil {
		log.Log(logger.ERROR, "Failed to set CPU frequency")
	}

	freq, err := machine.GetCPUFrequency()
	if err != nil {
		log.Log(logger.ERROR, "Failed to read CPU frequency")
	}

	println("[INFO] CPU frequency after setting:", freq)

	for {
		time.Sleep(time.Second)
		println("[INFO] MallardOS alive")
	}
}
