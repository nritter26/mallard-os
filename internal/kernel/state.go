package kernel

type State uint8

const (
	StateUninitialized State = iota
	StateInitializing
	StateReady
	StateStarting
	StateRunning
	StateStopping
	StateStopped
	StateFailed
)

func (s State) String() string {
	switch s {
		case StateUninitialized:
			return "uninitialized"
		case StateInitializing:
			return "initializing"
		case StateReady:
			return "ready"
		case StateStarting:
			return "starting"
		case StateRunning:
			return "running"
		case StateStopping:
			return "stopping"
		case StateStopped:
			return "stopped"
		case StateFailed:
			return "failed"
		default:
			return "unknown"
	}
}
