package kernel

type State uint8

const (
	StateUninitialized State = iota
	StateInitializing
	StateReady
	StateRunning
	StateStopped
	StateFailed
)

func (s State) String() string {
	switch s {

		case StateUninitialized:
			return "unitialized"
		case StateInitializing:
			return "initializing"
		case StateReady:
			return "ready"
		case StateRunning:
			return "running"
		case StateStopped:
			return "stopped"
		case StateFailed:
			return "failed"
		default:
			return "unknown"
	}
}

