package device

type State string

const (
	StateUninitialized State = "uninitialized"
	StateInitializing State = "initializing"
	StateReady State = "ready"
	StateStopping State = "stopping"
	StateStopped State = "stopped"
	StateFailed State = "failed"
)

