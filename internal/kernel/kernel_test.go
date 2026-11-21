package kernel

import (
	"errors"
	"testing"
	"strings"
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

func TestSubsystemInitializationOrder(t *testing.T) {
	k := New()

	var order []string

	subsystems := []string{
		"console",
		"storage",
		"network",
	}

	for _, name := range subsystems {
		subsystem := &mockSubsystem{
			name:  name,
			order: &order,
		}

		if err := k.Register(subsystem); err != nil {
			t.Fatal(err)
		}
	}

	if err := k.Initialize(); err != nil {
		t.Fatal(err)
	}

	if len(order) != len(subsystems) {
		t.Fatalf("initialized %d subsystems, want %d",
			 len(order), len(subsystems))
	}

	for i, name := range subsystems {
		if order[i] != name {
			t.Fatalf("initialization order[%d] = %s, want %s",
				 i, order[i], name)
		}
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

func TestInitializationStopsOnFailure(t *testing.T) {
	k := New()

	var order []string

	subsystems := []*mockSubsystem{
		{
			name:  "console",
			order: &order,
		},
		{
			name:   "storage",
			failed: true,
			order:  &order,
		},
		{
			name:  "network",
			order: &order,
		},
	}

	for _, subsystem := range subsystems {
		if err := k.Register(subsystem); err != nil {
			t.Fatal(err)
		}
	}

	err := k.Initialize()

	if err == nil {
		t.Fatal("Initialize() expected an error")
	}

	expected := []string{"console", "storage"}

	if len(order) != len(expected) {
		t.Fatalf("initialized %v, want %v", order, expected)
	}

	for i, name := range expected {
		if order[i] != name {
			t.Fatalf("initialization order[%d] = %s, want %s",
				 i, order[i], name)
		}
	}

	if subsystems[2].initialized {
		t.Fatal("network subsystem initialized after storage failure")
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
	order *[]string
	stopFailed bool
	stopCalled bool
	stopOrder *[]string
}

func (m *mockSubsystem) Name() string {
	return m.name
}

func (m *mockSubsystem) Init() error {
	m.initialized = true

	if m.order != nil {
		*m.order = append(*m.order, m.name)
	}

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

func (m *mockSubsystem) Stop() error {
	m.stopCalled = true

	if m.stopOrder != nil {
		*m.stopOrder = append(*m.stopOrder, m.name)
	}

	if m.stopFailed {
		return errors.New("subsystem shutdown failed")
	}

	return nil
}

func TestShutdownTransitionsToStopped(t *testing.T) {
	k := New()

	k.Register(&mockSubsystem{name: "console"})
	k.Initialize()
	k.Start()

	if err := k.Shutdown(); err != nil {
		t.Fatal(err)
	}

	if got := k.State(); got != StateStopped {
		t.Fatalf("State() = %v, want %v", got, StateStopped)
	}
}

func TestShutdownReverseOrder(t *testing.T) {
	k := New()

	var order []string

	names := []string{"console", "storage", "network"}

	for _, name := range names {
		k.Register(&mockSubsystem{
			name:      name,
			stopOrder: &order,
		})
	}

	k.Initialize()
	k.Start()

	if err := k.Shutdown(); err != nil {
		t.Fatal(err)
	}

	expected := []string{"network", "storage", "console"}

	for i, name := range expected {
		if order[i] != name {
			t.Fatalf("shutdown order[%d] = %s, want %s",
				 i, order[i], name)
		}
	}
}

func TestShutdownBeforeStart(t *testing.T) {
	k := New()
	k.Initialize()

	err := k.Shutdown()

	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("Shutdown() error = %v, want %v",
			 err, ErrInvalidState)
	}
}

func TestShutdownContinuesAfterFailure(t *testing.T) {
	k := New()

	var order []string

	subsystems := []*mockSubsystem{
		{name: "console", stopOrder: &order},
		{name: "storage", stopFailed: true, stopOrder: &order},
		{name: "network", stopOrder: &order},
	}

	for _, subsystem := range subsystems {
		k.Register(subsystem)
	}

	k.Initialize()
	k.Start()

	err := k.Shutdown()

	if err == nil {
		t.Fatal("Shutdown() expected an error")
	}

	if len(order) != 3 {
		t.Fatalf("stopped %d subsystems, want 3", len(order))
	}

	if got := k.State(); got != StateFailed {
		t.Fatalf("State() = %v, want %v", got, StateFailed)
	}
}

func TestShutdownPreservesMultipleErrors(t *testing.T) {
	k := New()

	k.Register(&mockSubsystem{
		name:       "console",
		stopFailed: true,
	})

	k.Register(&mockSubsystem{
		name:       "storage",
		stopFailed: true,
	})

	k.Initialize()
	k.Start()

	err := k.Shutdown()

	if err == nil {
		t.Fatal("Shutdown() expected an error")
	}

	if !strings.Contains(err.Error(), "console") {
		t.Error("shutdown error should contain console")
	}

	if !strings.Contains(err.Error(), "storage") {
		t.Error("shutdown error should contain storage")
	}
}
