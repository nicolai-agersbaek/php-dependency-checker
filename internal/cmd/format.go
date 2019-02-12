package cmd

import (
	"fmt"
	"time"
)

// FormatDuration formats the given duration, d.
func FormatDuration(d time.Duration) string {
	suffix := "ns"

	if d > time.Second {
		d = d / time.Second
		suffix = "s"
	} else if d > 10*time.Millisecond {
		d = d / time.Millisecond
		suffix = "ms"
	} else if d > 10*time.Microsecond {
		d = d / time.Microsecond
		suffix = "μs"
	}

	return fmt.Sprintf("%d"+suffix, d)
}
