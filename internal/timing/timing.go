package timing

import (
	"pizzago/internal/configs"
	"time"
)

func WaitFor(t uint) {
	time.Sleep(time.Duration(t) * time.Duration(configs.SlowFactor) * time.Millisecond)
}
