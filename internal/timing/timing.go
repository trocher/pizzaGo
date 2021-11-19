package timing

import (
	"pizzago/internal/configs"
	"time"
)

func WaitFor(t uint64) {
	time.Sleep(time.Duration(t) * time.Duration(configs.SlowFactor) * time.Millisecond)
}

func ExpectedTime() time.Duration {
	bottleneck := configs.Parameters.NumberOfWorker
	if bottleneck > configs.Parameters.NumberOfOven {
		bottleneck = configs.Parameters.NumberOfOven
	}
	expected := (configs.Timings.Bake + configs.Timings.Prepare + configs.Timings.Process + configs.Timings.QualityCheck) * configs.SlowFactor * configs.Parameters.NumberOfOrder / bottleneck
	return time.Duration(expected) * time.Duration(time.Millisecond)
}
