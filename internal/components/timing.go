package components

import (
	"time"
)

func WaitFor(t uint64) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}
