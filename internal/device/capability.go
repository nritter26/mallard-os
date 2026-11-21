package device

type Capability string

const (
	CapabilityGPIO Capability = "gpio"
	CapabilityUART Capability = "uart"
	CapabilityUSB Capability = "usb"
	CapabilitySPI Capability = "spi"
	CapabilityI2C Capability = "i2c"
	CapabilityPWM Capability = "pwm"
	CapabilityADC Capability = "adc"
	CapabilityDAC Capability = "dac"
	CapabilityWiFi Capability = "wifi"
	CapabilityBLE Capability = "ble"
	CapabilityPSRAM Capability = "psram"
	CapabilitySDCard Capability = "sdcard"
	CapabilityDisplay Capability = "display"
)

type CapabilityProvider interface {
	Capabilities() []Capability
}

