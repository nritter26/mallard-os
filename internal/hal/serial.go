package hal

import (
	"errors"
	"machine"
)

var (
	ErrConsoleNotInitialized = errors.New("console not initialized")
	ErrNoData = errors.New("no data available")
)

type USBConsole struct{
	initialized bool
}

var _ Console = (*USBConsole)(nil)

func NewUSBConsole() *USBConsole {
	return &USBConsole{}
}

func (s *USBConsole) Name() string {
	return "USB CDC"
}

func (s *USBConsole) Init() error {
	if err := machine.Serial.Configure(machine.UARTConfig{}); err != nil {
	    return err
	}
	s.initialized = true

	return nil
}

func (s *USBConsole) Write(data []byte) (int, error) {
	if !s.initialized {
		return 0, ErrConsoleNotInitialized
	}

	return machine.Serial.Write(data)
}

func (s *USBConsole) ReadByte() (byte, error) {
	if !s.initialized {
		return 0, ErrConsoleNotInitialized
	}

	return machine.Serial.ReadByte()
}

func (s *USBConsole) Buffered() int {
	if !s.initialized {
		return 0
	}

	return machine.Serial.Buffered()
}
