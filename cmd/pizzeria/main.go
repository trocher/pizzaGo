// This file run the pizzeria while timing it.
package main

import (
	"log"
	"pizzago/internal/components"
	"pizzago/internal/config"
	"time"
)

func main() {
	log.Printf("Starting pizzaGo ...")

	// Run the pizzeria while timing it before printing average throughput and latency
	start := time.Now()
	var cfg config.Config
	config.ReadConfig(&cfg)
	var latency = components.StartPizzeria(cfg)
	elapsed := time.Since(start)

	log.Printf("Took %s to cook %d pizzas with %d workers and %d ovens", elapsed, components.Config.Parameters.NumberOfOrders, components.Config.Parameters.NumberOfWorkers, components.Config.Parameters.NumberOfOvens)
	log.Printf("Average throuput of : %f pizza/ms and latency of %s ms/pizza", float64(components.Config.Parameters.NumberOfOrders)/float64(time.Duration(elapsed).Milliseconds()), latency)
}
