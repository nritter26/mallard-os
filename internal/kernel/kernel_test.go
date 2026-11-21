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

type mockSubsystem struct {
	name    string
	failed  bool
	started bool
}

func (m *mockSubsystem) Name() string {
	return m.name
}

func (m *mockSubsystem) Init() error {
	m.started = true

	if m.failed {
		return errors.New("subsystem initialization failed")
	}

	return nil
}
