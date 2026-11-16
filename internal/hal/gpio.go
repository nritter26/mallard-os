package hal

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
	name string
	pin GPIOPin
	mode PinMode
}

func NewGPIO(name string, pin GPIOPin, mode PinMode) *GPIODevice {
	return &GPIODevice{
		name: name,
		pin: pin,
		mode: mode,
	}
}

func (g *GPIODevice) Name() string {
	return g.name
}

func (g *GPIODevice) Init() error {
	g.pin.Configure(g.mode)
	return nil
}

func (g *GPIODevice) Write(value bool) {
	g.pin.Write(value)
}

func (g *GPIODevice) Read() bool {
	return g.pin.Read()
}



