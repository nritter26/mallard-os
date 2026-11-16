package hal

import "errors"

var (
	ErrNilPin        = errors.New("GPIO pin is nil")
	ErrInvalidMode   = errors.New("invalid GPIO mode")
	ErrNotInitialized = errors.New("GPIO not initialized")
	ErrInputOnly     = errors.New("cannot write to input GPIO")
)

type PinMode uint8

const (
	PinInput PinMode = iota
	PinOutput
)

type GPIOPin interface {
	Configure(mode PinMode)
	Write(value bool)
	Read() bool
}

type GPIODevice struct {
	name        string
	pin         GPIOPin
	mode        PinMode
	initialized bool
}

func NewGPIO(name string, pin GPIOPin, mode PinMode) *GPIODevice {
	return &GPIODevice{
		name: name,
		pin:  pin,
		mode: mode,
	}
}

func (g *GPIODevice) Name() string {
	return g.name
}

func (g *GPIODevice) Init() error {
	if g.pin == nil {
		return ErrNilPin
	}

	if g.mode != PinInput && g.mode != PinOutput {
		return ErrInvalidMode
	}

	g.pin.Configure(g.mode)
	g.initialized = true

	return nil
}

func (g *GPIODevice) Write(value bool) error {
	if !g.initialized {
		return ErrNotInitialized
	}

	if g.mode != PinOutput {
		return ErrInputOnly
	}

	g.pin.Write(value)
	return nil
}

func (g *GPIODevice) Read() (bool, error) {
	if !g.initialized {
		return false, ErrNotInitialized
	}

	return g.pin.Read(), nil
}
