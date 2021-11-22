package components

import (
	"time"
)

// Wait for t milliseconds
func WaitFor(t uint64) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}
