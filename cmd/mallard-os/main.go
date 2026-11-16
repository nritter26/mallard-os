package main

import (
    "mallard-os/internal/boot"
    "mallard-os/internal/hal"
)

func main() {
    console := hal.NewUSBConsole()

    manager := boot.NewManager(
        console,
    )

    if err := manager.Run(); err != nil {
        println("[FATAL] MallardOS boot failed")

        for {
        }
    }

    for {
    }
}
