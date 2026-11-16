package boot

import "mallard-os/internal/logger"

type BootStep interface {
	Name() string
	Init() error
}

type Manager struct {
	Steps  []BootStep
	Logger *logger.Logger
}

func NewManager(steps ...BootStep) *Manager {
	return &Manager{
		Steps: steps,
	}
}

func (m *Manager) SetLogger(log *logger.Logger) {
	m.Logger = log
}

func (m *Manager) log(level logger.Level, message string) {
	if m.Logger != nil {
		m.Logger.Log(level, message)
	}
}

func (m *Manager) Run() error {
	m.log(logger.INFO, "Starting MallardOS...")

	for _, step := range m.Steps {
		m.log(logger.INFO, "Initializing: "+step.Name())

		if err := step.Init(); err != nil {
			m.log(logger.ERROR, step.Name()+" initialization failed")
			return err
		}

		m.log(logger.INFO, step.Name()+" initialized")
	}

	m.log(logger.INFO, "Boot sequence completed.")
	return nil
}
