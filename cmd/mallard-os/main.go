package main

import (
	"time"

	"mallard-os/internal/boot"
	"mallard-os/internal/hal"
	"mallard-os/internal/logger"
)

func main() {
	time.Sleep(3 * time.Second)

	log := logger.New()

	console := hal.NewUSBConsole()

	hardware := hal.NewManager(console)

	manager := boot.NewManager(hardware)
	manager.SetLogger(log)

	if err := manager.Run(); err != nil {
		log.Log(logger.ERROR, "MallardOS boot failed")

		for {
		}
	}

	for {
	}
}

