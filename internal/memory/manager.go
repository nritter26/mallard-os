package memory

import (
	"errors"
	"sync"
)

var (
	ErrEmptyRegionName = errors.New("memory region name cannot be empty")
	ErrDuplicateRegion = errors.New("memory region already registered")
	ErrInvalidRegionCapacity = errors.New("available memory exceeds total capacity")
)

type Manager struct {
	mu sync.RWMutex
	regions map[string]Region
	order []string
}

func NewManager() *Manager {
	return &Manager{
		regions: make(map[string]Region),
	}
}

func (m *Manager) Register(region Region) error {
	if region.Name == "" {
		return ErrEmptyRegionName
	}

	if region.Available > region.Total {
		return ErrInvalidRegionCapacity
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.regions[region.Name]; exists {
		return ErrDuplicateRegion
	}

	m.regions[region.Name] = region
	m.order = append(m.order, region.Name)

	return nil
}

func (m *Manager) Get(name string) (Region, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	region, exists := m.regions[name]
	return region, exists
}

func (m *Manager) List() []Region {
	m.mu.RLock()
	defer m.mu.RUnlock()

	regions := make([]Region, 0, len(m.order))

	for _, name := range m.order {
		regions = append(regions, m.regions[name])
	}

	return regions
}

func (m *Manager) Name() string {
	return "Memory Manager"
}

func (m *Manager) Init() error {
	return nil
}

func (m *Manager) Stop() error {
	return nil
}
