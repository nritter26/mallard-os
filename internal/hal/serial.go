package hal

import "machine"

type USBConsole struct{}

func NewUSBConsole() *USBConsole {
    return &USBConsole{}
}

func (s *USBConsole) Name() string {
    return "USB CDC"
}

func (s *USBConsole) Init() error {
    return machine.Serial.Configure(machine.UARTConfig{})
}
