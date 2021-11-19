package main

import (
	"fmt"
	"log"
	"pizzago/internal/components"
	"pizzago/internal/configs"
	"strconv"
	"sync"
	"time"
)

func main() {
	fmt.Println("Starting pizzago")

	pizzaWorkers := make([]components.PizzaWorker, configs.Parameters.NumberOfWorker)
	fmt.Println("There are " + strconv.Itoa(configs.Parameters.NumberOfWorker) + " workers today")

	components.InitOven()
	var wg sync.WaitGroup

	for i := range pizzaWorkers {
		fmt.Println("Waking up worker " + strconv.Itoa(i))
		pizzaWorkers[i] = components.PizzaWorker{Name: uint64(i + 1), HasAssignedOven: false}
		wg.Add(1)
	}
	log.Printf("Starting the pizzeria")
	start := time.Now()

	for _, pizzaWorker := range pizzaWorkers {
		go pizzaWorker.Work(&wg)
	}
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Took %s to cook %d pizzas with %d workers and %d ovens", elapsed, configs.Parameters.NumberOfOrder, configs.Parameters.NumberOfWorker, configs.Parameters.NumberOfOven)
	log.Printf("%d order were taken and %d were delivered", components.OrderTaken, components.OrderDelivered)
}
