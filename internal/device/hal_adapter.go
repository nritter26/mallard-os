//go:build tinygo
package device

import "mallard-os/internal/hal"

type GPIOAdapter struct {
	gpio *hal.GPIODevice
	required bool
}

func NewGPIOAdapter(gpio *hal.GPIODevice, required bool) *GPIOAdapter {
	return &GPIOAdapter{
		gpio: gpio,
		required: required,
	}
}

func (a *GPIOAdapter) Name() string {
	return a.gpio.Name()
}

func (a *GPIOAdapter) Type() Type {
	return TypeGPIO
}

func (a *GPIOAdapter) Init() error {
	return a.gpio.Init()
}

func (a *GPIOAdapter) Required() bool {
	return a.required
}

func (a *GPIOAdapter) Capabilities() []Capability {
	return []Capability{
		CapabilityGPIO,
	}
}

type USBConsoleAdapter struct {
	console *hal.USBConsole
	required bool
}

func NewUSBConsoleAdapter(console *hal.USBConsole, required bool) *USBConsoleAdapter {
	return &USBConsoleAdapter{
		console: console,
		required: required,
	}
}

func (a *USBConsoleAdapter) Name() string {
	return a.console.Name()
}

func (a *USBConsoleAdapter) Type() Type {
	return TypeSerial
}

func (a *USBConsoleAdapter) Init() error {
	return a.console.Init()
}

func (a *USBConsoleAdapter) Required() bool {
	return a.required
}

func (a *USBConsoleAdapter) Capabilities() []Capability {
	return []Capability{
		CapabilityUART,
		CapabilityUSB,
	}
}
