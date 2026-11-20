package runtime

import (
	tinygoruntime "runtime"
	"time"
)

var bootTime time.Time

// Initialize records the Mallard-OS startup time
func Initialize() {
	bootTime = time.Now()
}

// Uptime returns the elapsed time since system startup
func Uptime() time.Duration {
	if bootTime.IsZero() {
		return 0
	}

	return time.Since(bootTime)
}

// MemoryStats returns the current TinyGo runtime memory statistics
func MemoryStats() tinygoruntime.MemStats {
	var stats tinygoruntime.MemStats
	tinygoruntime.ReadMemStats(&stats)
	return stats
}
