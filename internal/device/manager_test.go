package device

import (
	"errors"
	"testing"
)

type mockDevice struct {
	name string
	deviceType Type
	required bool
}

func(d *mockDevice) Name() string {return d.name}
func(d *mockDevice) Type() Type {return (d.deviceType)}
func(d *mockDevice) Required() bool {return (d.required)}
func(d *mockDevice) Init() error {return nil}

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
