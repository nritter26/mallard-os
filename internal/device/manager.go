package device

import (
	"sync"
	"errors"
)

var (
	ErrNilDevice = errors.New("device cannot be nil")
	ErrEmptyDeviceName = errors.New("device name cannot be empty")
	ErrDuplicatedDevice = errors.New("device already registered")

)

type Manager struct {
	mu      sync.RWMutex
	devices map[string]Device
	states map[string]State
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
		return ErrDuplicatedDevice
	}

	m.devices[name] = dev
	m.states[name] = StateUnintitialized

	return nil
}

