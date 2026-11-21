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

type mockStoppableDevice struct {
	*mockDevice
	stopErr error
	stopCalls int
	stopOrder *[]string
}

func (d *mockStoppableDevice) Stop() error {
	d.stopCalls++

	if d.stopOrder != nil {
		*d.stopOrder = append(*d.stopOrder, d.name)
	}

	return d.stopErr
}

func TestShutdownReverseOrder(t *testing.T) {
	m := NewManager()
	var order []string

	first := &mockStoppableDevice{
		mockDevice: &mockDevice{name: "console"},
		stopOrder:  &order,
	}

	second := &mockStoppableDevice{
		mockDevice: &mockDevice{name: "gpio"},
		stopOrder:  &order,
	}

	_ = m.Register(first)
	_ = m.Register(second)

	if err := m.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err := m.Shutdown(); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	expected := []string{"gpio", "console"}

	for i, name := range expected {
		if order[i] != name {
			t.Fatalf("shutdown order = %v, want %v", order, expected)
		}
	}
}

func TestShutdownContinuesAfterFailure(t *testing.T) {
	m := NewManager()
	failure := errors.New("shutdown failed")

	first := &mockStoppableDevice{
		mockDevice: &mockDevice{name: "console"},
		stopErr:    failure,
	}

	second := &mockStoppableDevice{
		mockDevice: &mockDevice{name: "gpio"},
	}

	_ = m.Register(first)
	_ = m.Register(second)

	if err := m.Initialize(); err != nil {
		t.Fatal(err)
	}

	err := m.Shutdown()

	if !errors.Is(err, failure) {
		t.Fatalf("Shutdown() error = %v, want %v", err, failure)
	}

	if first.stopCalls != 1 || second.stopCalls != 1 {
		t.Fatal("shutdown should attempt to stop every initialized device")
	}

	if m.states["console"] != StateFailed {
		t.Fatal("console should be marked failed")
	}

	if m.states["gpio"] != StateStopped {
		t.Fatal("gpio should be stopped")
	}
}

func TestShutdownDeviceWithoutStopper(t *testing.T) {
	m := NewManager()

	dev := &mockDevice{name: "gpio"}

	_ = m.Register(dev)

	if err := m.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err := m.Shutdown(); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	if m.states["gpio"] != StateStopped {
		t.Fatalf("state = %v, want %v", m.states["gpio"], StateStopped)
	}
}

func TestInitializeTwice(t *testing.T) {
	m := NewManager()
	dev := &mockDevice{name: "console"}

	_ = m.Register(dev)

	if err := m.Initialize(); err != nil {
		t.Fatal(err)
	}

	err := m.Initialize()

	if !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf("error = %v, want %v", err, ErrAlreadyInitialized)
	}

	if dev.initCalls != 1 {
		t.Fatalf("Init() called %d times, want 1", dev.initCalls)
	}
}

func TestRegisterAfterInitialize(t *testing.T) {
	m := NewManager()

	if err := m.Initialize(); err != nil {
		t.Fatal(err)
	}

	err := m.Register(&mockDevice{name: "gpio"})

	if !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf("error = %v, want %v", err, ErrAlreadyInitialized)
	}
}

func TestShutdownTwice(t *testing.T) {
	m := NewManager()

	dev := &mockStoppableDevice{
		mockDevice: &mockDevice{name: "console"},
	}

	_ = m.Register(dev)

	if err := m.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err := m.Shutdown(); err != nil {
		t.Fatal(err)
	}

	err := m.Shutdown()

	if !errors.Is(err, ErrManagerShuttingDown) {
		t.Fatalf("error = %v, want %v", err, ErrManagerShuttingDown)
	}

	if dev.stopCalls != 1 {
		t.Fatalf("Stop() called %d times, want 1", dev.stopCalls)
	}
}

type blockingDevice struct {
	*mockDevice
	started chan struct{}
	release chan struct{}
}

func (d *blockingDevice) Init() error {
	close(d.started)
	<-d.release
	return nil
}

func TestConcurrentInitialize(t *testing.T) {
	m := NewManager()

	dev := &blockingDevice{
		mockDevice: &mockDevice{
			name: "console",
			required: true,
		},
		started: make(chan struct{}),
		release: make(chan struct{}),
	}

	if err := m.Register(dev); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		done <-m.Initialize()
	}()

	<-dev.started

	err := m.Initialize()

	if !errors.Is(err, ErrInitializationInProgress) {
		t.Fatalf("error = %v, want %v,", err, ErrInitializationInProgress)
	}

	close(dev.release)

	if err := <-done; err != nil {
		t.Fatalf("first Initialize() error %v", err)
	}
}

func TestShutdownDuringInitialization(t *testing.T) {
	m := NewManager()

	dev := &blockingDevice{
		mockDevice: &mockDevice{
			name: "console",
			required: true,
		},
		started: make(chan struct{}),
		release: make(chan struct{}),
	}

	if err := m.Register(dev); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)

	go func() {
		done <- m.Initialize()
	}()

	<-dev.started

	err := m.Shutdown()

	if !errors.Is(err, ErrInitializationInProgress) {
		t.Fatalf("error = %v, want %v", err, ErrInitializationInProgress)
	}

	close(dev.release)

	if err := <-done; err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if m.states["console"] != StateReady {
		t.Fatal("device should be ready after initialization")
	}
}

func TestGetDevice(t *testing.T) {
	m := NewManager()

	dev := &mockDevice{
		name: "console",
		deviceType: TypeSerial,
		required: true,
	}

	if err := m.Register(dev); err != nil {
		t.Fatal(err)
	}

	got, state, exists := m.Get("console")

	if !exists {
		t.Fatal("device should exist")
	}

	if got.Name() != "console" {
		t.Fatalf("name = %q, want console", got.Name())
	}

	if state != StateUninitialized {
		t.Fatalf("state = %v, want %v", state, StateUninitialized)
	}

	_, _, exists = m.Get("unknown")

	if exists {
		t.Fatal("unknown device should not exist")
	}
}

func TestListDevices(t *testing.T) {
	m := NewManager()

	_ = m.Register(&mockDevice{
		name:       "console",
		deviceType: TypeSerial,
		required:   true,
	})

	_ = m.Register(&mockDevice{
		name:       "gpio",
		deviceType: TypeGPIO,
	})

	devices := m.List()

	if len(devices) != 2 {
		t.Fatalf("device count = %d, want 2", len(devices))
	}

	if devices[0].Name != "console" {
		t.Fatalf("first device = %q, want console", devices[0].Name)
	}

	if devices[1].Name != "gpio" {
		t.Fatalf("second device = %q, want gpio", devices[1].Name)
	}

	if !devices[0].Required {
		t.Fatal("console should be required")
	}
}

func TestGetDeviceState(t *testing.T) {
	m := NewManager()

	_ = m.Register(&mockDevice{name: "console"})

	state, exists := m.State("console")

	if !exists {
		t.Fatal("console should exist")
	}

	if state != StateUninitialized {
		t.Fatalf("state = %v, want %v", state, StateUninitialized)
	}

	_, exists = m.State("unknown")

	if exists {
		t.Fatal("unknown device should not exist")
	}
}

type capabilityTestDevice struct {
	name string
	capabilities []Capability
}

func(d *capabilityTestDevice) Name() string { return d.name }
func(d *capabilityTestDevice) Type() Type { return TypeGPIO }
func(d *capabilityTestDevice) Init() error { return nil }
func(d *capabilityTestDevice) Required() bool { return false }

func (d *capabilityTestDevice) Capabilities() []Capability {
	return d.capabilities
}

func TestHasCapability(t *testing.T) {
	manager := NewManager()

	dev := &capabilityTestDevice{
		name: "gpio0",
		capabilities: []Capability{
			CapabilityGPIO,
			CapabilityPWM,
		},
	}

	if err := manager.Register(dev); err != nil {
		t.Fatal(err)
	}

	if !manager.HasCapability(CapabilityGPIO) {
		t.Error("expected GPIO capability")
	}

	if !manager.HasCapability(CapabilityPWM) {
		t.Error("expected PWM capability")
	}

	if manager.HasCapability(CapabilityWiFi) {
		t.Error("unexpected WiFi capability")
	}
}

func TestFindByCapability(t *testing.T) {
	manager := NewManager()

	devices := []*capabilityTestDevice{
		{name: "gpio0", capabilities: []Capability{CapabilityGPIO, CapabilityPWM}},
		{name: "gpio1", capabilities: []Capability{CapabilityGPIO}},
		{name: "wifi0", capabilities: []Capability{CapabilityWiFi}},
	}

	for _, dev := range devices {
		if err := manager.Register(dev); err != nil {
			t.Fatal(err)
		}
	}

	found := manager.FindByCapability(CapabilityGPIO)

	if len(found) != 2 {
		t.Fatalf("expected 2 GPIO devices, got %d", len(found))
	}

	if found[0].Name() != "gpio0" || found[1].Name() != "gpio1" {
		t.Error("devices were not returned in registration order")
	}
}
