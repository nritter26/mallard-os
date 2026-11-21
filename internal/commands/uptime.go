package commands

import (
	"strconv"

	mallardruntime "mallard-os/internal/runtime"
	"mallard-os/internal/shell"
)

func Uptime(io shell.CommandIO, args []string) error {
	elapsed := mallardruntime.Uptime()

	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60

	if err := io.Println("System Uptime:"); err != nil {
		return err
	}

	if err := io.Println(
		"  Running: " +
		formatDuration(hours, minutes, seconds),
	); err != nil {
		return err
	}

	return io.Println(
		"  Seconds: " + strconv.Itoa(int(elapsed.Seconds())),
	)
}

func formatDuration(hours, minutes, seconds int) string {
	return twoDigits(hours) + ":" +
	twoDigits(minutes) + ":" +
	twoDigits(seconds)
}

func twoDigits(value int) string {
	if value < 10 {
		return "0" + strconv.Itoa(value)
	}

	return strconv.Itoa(value)
}
