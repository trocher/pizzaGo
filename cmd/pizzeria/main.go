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

	// Run the pizzeria while timing it to be able to benchmark its performances
	start := time.Now()
	var cfg config.Config
	config.ReadConfig(&cfg)
	components.StartPizzeria(cfg)
	elapsed := time.Since(start)
	expectedTime := components.ExpectedTime()

	log.Printf("Took %s to cook %d pizzas with %d workers and %d ovens", elapsed, components.Config.Parameters.NumberOfOrders, components.Config.Parameters.NumberOfWorkers, components.Config.Parameters.NumberOfOvens)
	log.Printf("overhead was %s (%d%%) as time taken is %s and expected time would be %s ", elapsed-expectedTime, int((float64(elapsed)/float64(expectedTime))*100)-100, elapsed, expectedTime)
}
