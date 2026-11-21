package kernel

import (
	"errors"
	"testing"
)

func TestNewKernelStartsUninitialized(t *testing.T) {
	k := New()

	if got := k.State(); got != StateUninitialized {
		t.Fatalf("State() = %v, want %v", got, StateUninitialized)
	}
}

func TestAlreadyInitialized(t *testing.T) {
	k := New()
	if err := k.Initialize(); err != nil {
		t.Fatalf("Initialize() returned unexpected error: %v", err)
	}
	err := k.Initialize()
	if !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf("Initialize() error = %v, want %v", err, ErrAlreadyInitialized)
	}
}

func TestInitializeTransitionsToReady(t *testing.T) {
	k := New()

	if err := k.Initialize(); err != nil {
		t.Fatalf("Initialized() returned unexpected error: %v", err)
	}

	if got := k.State(); got != StateReady {
		t.Fatalf("State() = %v, want %v", got, StateReady)
	}
}

func TestStartTransitionsToRunning(t *testing.T) {
	k := New()

	if err := k.Initialize(); err != nil {
		t.Fatalf("Initialized() returned unexpected error: %v", err)
	}

	if err := k.Start(); err != nil {
		t.Fatalf("Start() returned unexpected error: %v", err)
	}

	if got := k.State(); got != StateRunning {
		t.Fatalf("State() = %v, want %v", got, StateRunning)
	}
}

func TestStartBeforeInitializedReturnsErrInvalidState(t *testing.T) {
	k := New()

	err := k.Start()
	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("Start() error = %v, want %v", err, ErrInvalidState)
	}

	// A failed Start must not change state
	if got := k.State(); got != StateUninitialized {
		t.Fatalf("State() = %v, want %v", got, StateUninitialized)
	}
}

func TestRegisterSubsystem(t *testing.T) {
	k := New()
	subsystem := &mockSubsystem{name: "storage"}

	if err := k.Register(subsystem); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
}

func TestRegisterDuplicateSubsystem(t *testing.T) {
	k := New()

	subsystem := &mockSubsystem{name: "storage"}

	k.Register(subsystem)

	err := k.Register(&mockSubsystem{name: "storage"})

	if !errors.Is(err, ErrDuplicateSubsystem) {
		t.Fatalf("Register() error = %v, want %v",
			 err, ErrDuplicateSubsystem)
	}
}

func TestRegisterNilSubsystem(t *testing.T) {
	k := New()

	if err := k.Register(nil); !errors.Is(err, ErrInvalidSubsystem) {
		t.Fatalf("Register() error = %v, want %v",
			 err, ErrInvalidSubsystem)
	}
}

func TestRegisterAfterInitialization(t *testing.T) {
	k := New()
	k.Initialize()

	err := k.Register(&mockSubsystem{name: "storage"})

	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("Register() error = %v, want %v",
			 err, ErrInvalidState)
	}
}

func TestInitializationFailure(t *testing.T) {
	k := New()

	subsystem := &mockSubsystem{
		name:   "storage",
		failed: true,
	}

	k.Register(subsystem)

	err := k.Initialize()

	if err == nil {
		t.Fatal("Initialize() expected an error")
	}

	if got := k.State(); got != StateFailed {
		t.Fatalf("State() = %v, want %v", got, StateFailed)
	}
}

func TestSubsystemStart(t *testing.T) {
	k := New()

	subsystem := &mockSubsystem{name: "storage"}
	k.Register(subsystem)

	if err := k.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err := k.Start(); err != nil {
		t.Fatal(err)
	}

	if !subsystem.startCalled {
		t.Fatal("subsystem Start() was not called")
	}
}

func TestSubsystemStartFailure(t *testing.T) {
	k := New()

	subsystem := &mockSubsystem{
		name:        "storage",
		startFailed: true,
	}

	k.Register(subsystem)
	k.Initialize()

	err := k.Start()

	if err == nil {
		t.Fatal("Start() expected an error")
	}

	if got := k.State(); got != StateFailed {
		t.Fatalf("State() = %v, want %v", got, StateFailed)
	}
}

type mockSubsystem struct {
	name        string
	failed      bool
	startFailed bool
	initialized bool
	startCalled bool
}

func (m *mockSubsystem) Name() string {
	return m.name
}

func (m *mockSubsystem) Init() error {
	m.initialized = true

	if m.failed {
		return errors.New("subsystem initialization failed")
	}

	return nil
}

func (m *mockSubsystem) Start() error {
	m.startCalled = true

	if m.startFailed {
		return errors.New("subsystem startup failed")
	}

	return nil
}
