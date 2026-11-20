package console

import "errors"

var (
	ErrNotInitialized = errors.New("console transport not initialized")
)

// Transport defines the hardware communication capabilities
// required by the console layer
// HAL implementaions satisfy this interface

type Transport interface {
	ReadByte() (byte, error)
	Write(data []byte) (int, error)
	Buffered() int
}

type Console struct {
	transport Transport
}

func New(transport Transport) *Console {
	return &Console{
		transport: transport,
	}
}

func (c *Console) Write(data []byte) (int, error) {
	if c.transport == nil {
		return 0, ErrNotInitialized
	}
	return c.transport.Write(data)
}

func (c *Console) ReadByte() (byte, error) {
	if c.transport == nil {
		return 0, ErrNotInitialized
	}
	return c.transport.ReadByte()
}

func (c *Console) Buffered() int {
	if c.transport == nil {
		return 0
	}

	return c.transport.Buffered()
}

func (c *Console) Print(message string) error {
	_, err := c.Write([]byte(message))
	return err
}

func (c *Console) Println(message string) error {
	return c.Print(message + "\r\n")
}

