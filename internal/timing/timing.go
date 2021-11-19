package timing

import "time"

func WaitFor(t uint) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}
