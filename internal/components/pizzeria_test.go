package components

import (
	"fmt"
	"pizzago/internal/config"
	"testing"
)

func BenchmarkPizzeria(b *testing.B) {
	var cfg = config.Config{}

	// Set the config parameters
	cfg.Times.Process = uint64(1)
	cfg.Times.Prepare = uint64(2)
	cfg.Times.Bake = uint64(5)
	cfg.Times.QualityCheck = uint64(1)
	cfg.Parameters.NumberOfOrders = uint64(200)

	// First variable of the benchmark
	for nbWorker := uint64(1); nbWorker <= 10; nbWorker++ {
		cfg.Parameters.NumberOfWorkers = nbWorker

		// Second variable of the benchmark
		for nbOven := uint64(1); nbOven <= 10; nbOven++ {
			cfg.Parameters.NumberOfOvens = nbOven
			b.Run(fmt.Sprintf("%d_Workers;%d_Ovens", nbWorker, nbOven), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					StartPizzeria(cfg)
				}
			})
		}
	}
}
