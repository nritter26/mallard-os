package hal

import "fmt"

type Device interface {
	Name() string
	Init() error
}

type Manager struct {
	Devices []Device
}

func NewManager(devices ...Device) *Manager {
	return &Manager {
		Devices: devices,
	}
}

func (m *Manager) Init() error {
	for _, device := range m.Devices {
		if err := device.Init(); err != nil {
			return fmt.Errorf(
				"hal: device %s initialization failed: %w",
				device.Name(),
				err,
			)
		}
	}

	return  nil
}

func (m *Manager) Name() string {
	return "Hardware Abstraction Layer"
}
