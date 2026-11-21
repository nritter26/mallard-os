package kernel

type Subsystem interface {
	Name() string
	Init() error
}

// Starter is implemented by subsystems that need an activation phase
type Starter interface {
	Start() error
}

// Stopper is implemented by subsystems that support graceful shutdown
type Stopper interface {
	Stop() error
}
