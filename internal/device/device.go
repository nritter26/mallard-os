package device

type Type string

const (
	TypeGPIO Type = "gpio"
	TypeSerial Type = "serial"
	TypeSPI Type = "spi"
	TypeI2C Type = "i2c"
	TypePWM Type = "pwm"
	TypeADC Type = "adc"
)

type Device interface {
	Name() string
	Type() Type
	Init() error
	Required() bool
}

type Stopper interface {
	Stop() error
}

