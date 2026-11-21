package kernel

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyInitialized = errors.New("kernel already initialized")
	ErrInvalidState       = errors.New("invalid kernel state")
	ErrDuplicateSubsystem = errors.New("subsystem already registered")
	ErrInvalidSubsystem   = errors.New("invalid subsystem")
)

type Kernel struct {
	mu         sync.Mutex
	state      State
	subsystems []Subsystem
}

func New() *Kernel {
	return &Kernel{
		state: StateUninitialized,
	}
}

func (k *Kernel) State() State {
	k.mu.Lock()
	defer k.mu.Unlock()

	return k.state
}

// Register adds a subsystem to the kernel
func (k *Kernel) Register(subsystem Subsystem) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateUninitialized {
		return ErrInvalidState
	}

	if subsystem == nil || subsystem.Name() == "" {
		return ErrInvalidSubsystem
	}

	for _, existing := range k.subsystems {
		if existing.Name() == subsystem.Name() {
			return ErrDuplicateSubsystem
		}
	}

	k.subsystems = append(k.subsystems, subsystem)

	return nil
}

// Initialize initializes all registered subsystems
func (k *Kernel) Initialize() error {
	k.mu.Lock()

	if k.state != StateUninitialized {
		k.mu.Unlock()
		return ErrAlreadyInitialized
	}

	k.state = StateInitializing
	subsystems := append([]Subsystem(nil), k.subsystems...)

	k.mu.Unlock()

	for _, subsystem := range subsystems {
		if err := subsystem.Init(); err != nil {
			k.mu.Lock()
			k.state = StateFailed
			k.mu.Unlock()

			return err
		}
	}

	k.mu.Lock()
	k.state = StateReady
	k.mu.Unlock()

	return nil
}

// Start transitions the kernel into running state
func (k *Kernel) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateReady {
		return ErrInvalidState
	}

	k.state = StateRunning

	return nil
}
