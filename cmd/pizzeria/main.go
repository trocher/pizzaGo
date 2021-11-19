package main

import (
	"log"
	"pizzago/internal/components"
	"time"
)

func main() {
	log.Printf("Starting pizzago ...")

	start := time.Now()
	components.StartPizzeria()
	elapsed := time.Since(start)

	log.Printf("Took %s to cook %d pizzas with %d workers and %d ovens", elapsed, components.Config.Parameters.NumberOfOrders, components.Config.Parameters.NumberOfWorkers, components.Config.Parameters.NumberOfOvens)

	expectedTime := components.ExpectedTime()

	log.Printf("overhead was %s (%d%%) as time taken is %s and expected time would be %s ", elapsed-expectedTime, int((float64(elapsed)/float64(expectedTime))*100)-100, elapsed, expectedTime)
}
