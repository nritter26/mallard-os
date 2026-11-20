package cmd

import (
	"fmt"

	mallardruntime "mallard-os/internal/runtime"
	"mallard-os/internal/shell"
)

// Meminfo displays current runtime memory statistics
func Meminfo(io shell.CommandIO, args []string) error {
	stats := mallardruntime.MemoryStats()

	lines := []string{
		"",
		"Mallard-OS Memory Information",
		"-----------------------------",
		fmt.Sprintf("Heap Allocated: %d bytes", stats.HeapAlloc),
		fmt.Sprintf("Heap In Use:    %d bytes", stats.HeapInuse),
		fmt.Sprintf("Heap Idle:      %d bytes", stats.HeapIdle),
		fmt.Sprintf("Heap Total:     %d bytes", stats.HeapSys),
		fmt.Sprintf("Heap Objects:   %d", stats.HeapObjects),
		fmt.Sprintf("Total Alloc:    %d bytes", stats.TotalAlloc),
		fmt.Sprintf("GC Cycles:      %d", stats.NumGC),
		"-----------------------------",
	}

	for _, line := range lines {
		if err := io.Println(line); err != nil {
			return err
		}
	}

	return nil
}
