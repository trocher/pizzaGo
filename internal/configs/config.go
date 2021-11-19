package configs

type TimingsConfig struct {
	Process      uint64
	Prepare      uint64
	Bake         uint64
	QualityCheck uint64
}

var Timings = TimingsConfig{1, 2, 5, 1}

type HyperParameters struct {
	NumberOfWorker uint64
	NumberOfOven   uint64
	NumberOfOrder  uint64
}

var Parameters = HyperParameters{2, 2, 500}

var SlowFactor uint64 = 1
