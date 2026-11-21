package device

import (
	"sync"
	"errors"
	"fmt"
)

var (
	ErrNilDevice = errors.New("device cannot be nil")
	ErrEmptyDeviceName = errors.New("device name cannot be empty")
	ErrDuplicateDevice = errors.New("device already registered")

)

type Manager struct {
	mu      sync.RWMutex
	devices map[string]Device
	states map[string]State
	initialized []string
	order []string
}

func NewManager() *Manager {
	return &Manager{
		devices: make(map[string]Device),
		states: make(map[string]State),
	}
}

func (m *Manager) Register(dev Device) error {
	if dev == nil {
		return ErrNilDevice
	}

	name := dev.Name()

	if name == "" {
		return ErrEmptyDeviceName
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.devices[name]; exists {
		return ErrDuplicateDevice
	}

	m.devices[name] = dev
	m.states[name] = StateUninitialized

	m.order = append(m.order, name)

	return nil
}

func (m *Manager) Initialize() error {
	m.mu.Lock()

	order := append([]string(nil), m.order...)
	devices := make(map[string]Device, len(m.devices))

	for name, dev := range m.devices {
		devices[name] = dev
	}

	m.mu.Unlock()

	for _, name := range order {
		dev := devices[name]

		m.mu.Lock()
		m.states[name] = StateInitializing
		m.mu.Unlock()

		if err := dev.Init(); err != nil {
			m.mu.Lock()
			m.states[name] = StateFailed
			m.mu.Unlock()

			if dev.Required() {
				return fmt.Errorf(
					"initialize required device %q: %w",
		      name,
		      err,
				)
			}

			continue
		}

		m.mu.Lock()
		m.states[name] = StateReady
		m.initialized = append(m.initialized, name)
		m.mu.Unlock()
	}

	return nil
}
