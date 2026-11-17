package hal

type Console interface {
	ReadByte() (byte, error)
	Write(data []byte) (int, error)
	Buffered() int
}
