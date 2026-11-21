package kernel

type Subsystem interface {
	Name() string
	Init() error
}

