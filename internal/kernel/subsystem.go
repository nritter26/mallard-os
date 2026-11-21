package kernel

type Subsystem interface {
	Name() string
	Init() error
}

// Starter is implemented by subsystems that need an activation phase
type Starter interface {
	Start() error
}
