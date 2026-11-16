package boot

import (
	"errors"
	"testing"
)

type mockStep struct {
	name string
	failed bool
	started bool
}

func (m *mockStep) Name() string {
	return m.name
}


func (m *mockStep) Init() error {
	m.started = true

	if m.failed {
		return errors.New("initialization failed")
	}

	return nil
}

func TestBootManager(t *testing.T) {
	runtime := &mockStep{name: "Runtime"}
	hal := &mockStep{name: "HAL"}
	kernel := &mockStep{name: "Kernel"}

	manager := NewManager(runtime, hal, kernel)

	if err := manager.Run(); err != nil {
		t.Fatal(err)
	}

	if !runtime.started || !hal.started || !kernel.started {
		t.Fatal("not all boot steps excuted!!")
	}
}

func TestBootFailure(t *testing.T) {
	runtime := &mockStep{name: "Runtime"}
	hal := &mockStep{name: "HAL", failed: true}
	kernel := &mockStep{name: "Kernel"}
	
	manager := NewManager(runtime, hal, kernel)

	if err := manager.Run(); err == nil {
		t.Fatal("expected boot failure")
	}

	if kernel.started {
		t.Fatal("kernel stared despite HAL failure")
	}
}

