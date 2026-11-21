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
	ErrAlreadyInitialized = errors.New("device manager already initialized")
	ErrInitializationInProgress = errors.New("device initialization in progress")
	ErrManagerShuttingDown = errors.New("device manager is shutting down")
)

type Manager struct {
	mu      sync.RWMutex
	devices map[string]Device
	states map[string]State
	initialized []string
	order []string
	initializationStarted bool
	initializing bool
	shuttingDown bool
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

	if m.shuttingDown {
		return ErrManagerShuttingDown
	}

	if m.initializationStarted {
		return ErrAlreadyInitialized
	}

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

	if m.shuttingDown {
		m.mu.Unlock()
		return ErrManagerShuttingDown
	}

	if m.initializing {
		m.mu.Unlock()
		return ErrInitializationInProgress
	}

	if m.initializationStarted {
		m.mu.Unlock()
		return ErrAlreadyInitialized
	}

	m.initializationStarted = true
	m.initializing = true

	order := append([]string(nil), m.order...)

	devices := make(map[string]Device, len(m.devices))

	for name, dev := range m.devices {
		devices[name] = dev
	}
	defer func() {
		m.mu.Lock()
		m.initializing = false
		m.mu.Unlock()
	}()

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

func (m *Manager) Shutdown() error {
	m.mu.Lock()

	if m.initializing {
		m.mu.Unlock()
		return ErrInitializationInProgress
	}

	if m.shuttingDown {
		m.mu.Unlock()
		return ErrManagerShuttingDown
	}

	m.shuttingDown = true

	var devices []string

	for i := len(m.initialized) - 1; i >= 0; i-- {
		name := m.initialized[i]

		if m.states[name] != StateReady {
			continue
		}

		m.states[name] = StateStopping
		devices = append(devices, name)
	}

	m.mu.Unlock()

	var shutdownErrors []error

	for _, name := range devices {
		m.mu.RLock()
		dev := m.devices[name]
		m.mu.RUnlock()

		stopper, ok := dev.(Stopper)
		if ok {
			if err := stopper.Stop(); err != nil {
				shutdownErrors = append(
					shutdownErrors,
			    fmt.Errorf("stop device %q: %w", name, err),
				)

				m.mu.Lock()
				m.states[name] = StateFailed
				m.mu.Unlock()

				continue
			}
		}

		m.mu.Lock()
		m.states[name] = StateStopped
		m.mu.Unlock()
	}

	return errors.Join(shutdownErrors...)
}

func (m *Manager) Get(name string) (Device, State, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	dev, exists := m.devices[name]
	if !exists {
		return nil, "", false
	}

	return dev, m.states[name], true
}

func (m *Manager) List() []DeviceInfo {
	m.mu.RLock()

	type entry struct {
		dev   Device
		state State
	}

	entries := make([]entry, 0, len(m.order))

	for _, name := range m.order {
		entries = append(entries, entry{
			dev: m.devices[name],
			state: m.states[name],
		})
	}

	m.mu.RUnlock()

	result := make([]DeviceInfo, 0, len(entries))

	for _, entry := range entries {
		result = append(result, DeviceInfo{
			Name: entry.dev.Name(),
				Type: entry.dev.Type(),
				State: entry.state,
				Required: entry.dev.Required(),
		})
	}

	return result
}

func (m *Manager) State(name string) (State, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[name]
	return state, exists
}
