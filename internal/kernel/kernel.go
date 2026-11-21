package kernel

import (
	"errors"
	"sync"
	"fmt"
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

//Register adds a subsystem to the kernel
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

//Initialize initializes all registered subsystems
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

//Start activates registered subsystems
func (k *Kernel) Start() error {
	k.mu.Lock()

	if k.state != StateReady {
		k.mu.Unlock()
		return ErrInvalidState
	}

	k.state = StateStarting
	subsystems := append([]Subsystem(nil), k.subsystems...)

	k.mu.Unlock()

	for _, subsystem := range subsystems {
		starter, ok := subsystem.(Starter)
		if !ok {
			continue
		}

		if err := starter.Start(); err != nil {
			k.mu.Lock()
			k.state = StateFailed
			k.mu.Unlock()

			return fmt.Errorf(
				"failed to start subsystem %s: %w",
		     subsystem.Name(),
					  err,
			)
		}
	}

	k.mu.Lock()
	k.state = StateRunning
	k.mu.Unlock()

	return nil
}

//Shutdown gracefully stops all registered subsystems
func (k *Kernel) Shutdown() error {
	k.mu.Lock()

	if k.state != StateRunning {
		k.mu.Unlock()
		return ErrInvalidState
	}

	// Reserve shutdown before releasing the mutex
	k.state = StateStopping

	subsystems := append([]Subsystem(nil), k.subsystems...)
	k.mu.Unlock()

	var shutdownErr error

	for i := len(subsystems) - 1; i >= 0; i-- {
		stopper, ok := subsystems[i].(Stopper)
		if !ok {
			continue
		}

		if err := stopper.Stop(); err != nil {
			shutdownErr = errors.Join(
				shutdownErr,
			     fmt.Errorf(
				     "failed to stop subsystem %s: %w",
		   subsystems[i].Name(),
					err,
			     ),
			)
		}
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	if shutdownErr != nil {
		k.state = StateFailed
		return shutdownErr
	}

	k.state = StateStopped
	return nil
}
