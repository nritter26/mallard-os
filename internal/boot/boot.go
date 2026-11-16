package boot

type BootStep interface {
	Name() string
	Init() error
}

type Manager struct {
	Steps []BootStep
}

func NewManager(steps...BootStep) *Manager {
	return &Manager{
		Steps: steps,
	}
}

func (m *Manager) Run() error {
    println("[BOOT] Starting Mallard-OS...")

    for _, step := range m.Steps {
        println("[BOOT] Initializing:", step.Name())

        if err := step.Init(); err != nil {
            println("[FAIL]", step.Name())
            return err
        }

        println("[ OK ]", step.Name(), "initialized")
    }

    println("[INFO] Boot sequence completed.")

    return nil
}
