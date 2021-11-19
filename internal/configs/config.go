package configs

type TimingsConfig struct {
	Process      uint
	Prepare      uint
	Bake         uint
	QualityCheck uint
}

var Timings = TimingsConfig{1, 2, 5, 1}

type HyperParameters struct {
	NumberOfWorker int
	NumberOfOven   uint64
	NumberOfOrder  uint64
}

var Parameters = HyperParameters{2, 4, 500}
