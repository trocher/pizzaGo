package main

import (
	"log"
	"pizzago/internal/components"
	"pizzago/internal/configs"
	"pizzago/internal/timing"
	"time"
)

func main() {
	log.Printf("Starting pizzago ...")

	start := time.Now()
	components.StartPizzeria()
	elapsed := time.Since(start)

	log.Printf("Took %s to cook %d pizzas with %d workers and %d ovens", elapsed, configs.Parameters.NumberOfOrder, configs.Parameters.NumberOfWorker, configs.Parameters.NumberOfOven)

	expectedTime := timing.ExpectedTime()

	log.Printf("overhead was %s (%d%%)as time taken is %s and expected time would be %s ", elapsed-expectedTime, elapsed-expectedTime/expectedTime, elapsed, expectedTime)
}
