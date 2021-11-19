package components

import (
	"time"
)

func WaitFor(t uint64) {
	time.Sleep(time.Duration(t) * time.Millisecond)
}

func ExpectedTime() time.Duration {
	bottleneck := Config.Parameters.NumberOfWorkers
	if bottleneck > Config.Parameters.NumberOfOvens {
		bottleneck = Config.Parameters.NumberOfOvens
	}
	expected := (Config.Times.Bake + Config.Times.Prepare + Config.Times.Process + Config.Times.QualityCheck) * Config.Parameters.NumberOfOrders / bottleneck
	return time.Duration(expected) * time.Duration(time.Millisecond)
}
