package hal

import "machine"

type machineGPIOPin struct {
	pin machine.Pin
}

func (p machineGPIOPin) Configure(mode PinMode) {
	machineMode := machine.PinInput

	if mode == PinOutput {
		machineMode = machine.PinOutput
	}

	p.pin.Configure(machine.PinConfig{
		Mode: machineMode,
	})
}

func (p machineGPIOPin) Write(value bool) {
	p.pin.Set(value)
}

func (p machineGPIOPin) Read() bool {
	return p.pin.Get()
}
