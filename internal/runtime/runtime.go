package runtime

import "time"

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
