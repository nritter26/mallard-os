package device

import (
	"errors"
	"testing"
)

type mockDevice struct {
	name string
	deviceType Type
	required bool
	initErr error
	initCalls int
}

func(d *mockDevice) Name() string {return d.name}
func(d *mockDevice) Type() Type {return (d.deviceType)}
func(d *mockDevice) Required() bool {return (d.required)}
func(d *mockDevice) Init() error{
	d.initCalls++
	return d.initErr
}

func TestRegisterDevice(t *testing.T) {
	m := NewManager()

	dev := &mockDevice{
		name:       "console",
		deviceType: TypeSerial,
		required:   true,
	}

	if err := m.Register(dev); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if _, exists := m.devices["console"]; !exists {
		t.Fatal("device was not registered")
	}

	if got := m.states["console"]; got != StateUninitialized {
		t.Fatalf("state = %v, want %v", got, StateUninitialized)
	}
}

func TestRegisterNilDevice(t *testing.T) {
	m := NewManager()

	err := m.Register(nil)

	if !errors.Is(err, ErrNilDevice) {
		t.Fatalf("Register(nil) error = %v, want %v", err, ErrNilDevice)
	}
}

func TestRegisterEmptyName(t *testing.T) {
	m := NewManager()

	dev := &mockDevice{}

	err := m.Register(dev)

	if !errors.Is(err, ErrEmptyDeviceName) {
		t.Fatalf("Register() error = %v, want %v", err, ErrEmptyDeviceName)
	}
}

func TestRegisterDuplicateDevice(t *testing.T) {
	m := NewManager()

	dev := &mockDevice{name: "console"}

	if err := m.Register(dev); err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	err := m.Register(dev)

	if !errors.Is(err, ErrDuplicateDevice) {
		t.Fatalf("second Register() error = %v, want %v", err, ErrDuplicateDevice)
	}
}

func TestInitializeDevices(t *testing.T) {
	m := NewManager()

	first := &mockDevice{name: "console", required: true}
	second := &mockDevice{name: "gpio", required: false}

	_ = m.Register(first)
	_ = m.Register(second)

	if err := m.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if m.states["console"] != StateReady {
		t.Fatal("console should be ready")
	}

	if m.states["gpio"] != StateReady {
		t.Fatal("gpio should be ready")
	}

	if first.initCalls != 1 || second.initCalls != 1 {
		t.Fatal("all devices should initialize exactly once")
	}
}

func TestInitializeRequiredDeviceFailure(t *testing.T) {
	m := NewManager()

	failure := errors.New("hardware initialization failed")

	required := &mockDevice{
		name:     "console",
		required: true,
		initErr:  failure,
	}

	next := &mockDevice{name: "gpio"}

	_ = m.Register(required)
	_ = m.Register(next)

	err := m.Initialize()

	if !errors.Is(err, failure) {
		t.Fatalf("Initialize() error = %v, want %v", err, failure)
	}

	if m.states["console"] != StateFailed {
		t.Fatal("required device should be failed")
	}

	if next.initCalls != 0 {
		t.Fatal("initialization should stop after required device failure")
	}
}

func TestInitializeOptionalDeviceFailure(t *testing.T) {
	m := NewManager()

	failure := errors.New("optional device unavailable")

	optional := &mockDevice{
		name:    "spi",
		initErr: failure,
	}

	next := &mockDevice{name: "gpio"}

	_ = m.Register(optional)
	_ = m.Register(next)

	if err := m.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if m.states["spi"] != StateFailed {
		t.Fatal("optional device should be marked failed")
	}

	if m.states["gpio"] != StateReady {
		t.Fatal("manager should continue initializing other devices")
	}
}

func TestInitializeRegistrationOrder(t *testing.T) {
	m := NewManager()

	first := &mockDevice{name: "first"}
	second := &mockDevice{name: "second"}
	third := &mockDevice{name: "third"}

	_ = m.Register(first)
	_ = m.Register(second)
	_ = m.Register(third)

	if err := m.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if len(m.initialized) != 3 {
		t.Fatalf("initialized count = %d, want 3", len(m.initialized))
	}

	expected := []string{"first", "second", "third"}

	for i, name := range expected {
		if m.initialized[i] != name {
			t.Fatalf("initialization order = %v, want %v", m.initialized, expected)
		}
	}
}
