package memory

import (
	"errors"
	"testing"
)

func TestRegisterRegion(t *testing.T) {
	manager := NewManager()

	region := Region{
		Name: "internal-sram",
		Type: RegionInternalRAM,
		Total: 512 * 1024,
		Available: 256 * 1024,
	}

	if err := manager.Register(region); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	got, exists := manager.Get("internal-sram")
	if !exists {
		t.Fatal("registered region not found")
	}

	if got != region {
		t.Errorf("Get() = %+v, want %+v", got, region)
	}
}

func TestRegisterDuplicateRegion(t *testing.T) {
	manager := NewManager()

	region := Region{
		Name: "internal-sram",
		Type: RegionInternalRAM,
	}

	if err := manager.Register(region); err != nil {
		t.Fatal(err)
	}

	err := manager.Register(region)
	if !errors.Is(err, ErrDuplicateRegion) {
		t.Errorf("Register() error = %v, want %v", err, ErrDuplicateRegion)
	}
}

func TestRegisterEmptyRegionName(t *testing.T) {
	manager := NewManager()

	err := manager.Register(Region{})
	if !errors.Is(err, ErrEmptyRegionName) {
		t.Errorf("Register() error = %v, want %v", err, ErrEmptyRegionName)
	}
}

func TestListRegions(t *testing.T) {
	manager := NewManager()

	regions := []Region{
		{Name: "internal-sram", Type: RegionInternalRAM},
		{Name: "external-psram", Type: RegionExternalRAM},
		{Name: "flash", Type: RegionFlash},
	}

	for _, region := range regions {
		if err := manager.Register(region); err != nil {
			t.Fatal(err)
		}
	}

	got := manager.List()

	if len(got) != len(regions) {
		t.Fatalf("List() returned %d regions, want %d", len(got), len(regions))
	}

	for i, region := range regions {
		if got[i] != region {
			t.Errorf("List()[%d] = %+v, want %+v", i, got[i], region)
		}
	}
}

func TestRegisterInvalidRegionCapacity(t *testing.T) {
	manager := NewManager()

	err := manager.Register(Region{
		Name: "invalid-region",
		Type: RegionInternalRAM,
		Total: 1024,
		Available: 2048,
	})

	if !errors.Is(err, ErrInvalidRegionCapacity) {
		t.Errorf("Register() error = %v, want %v",
			 err, ErrInvalidRegionCapacity)
	}
}
