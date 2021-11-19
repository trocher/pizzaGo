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

	for i := 0; i < configs.Parameters.NumberOfWorker; i++ {
		fmt.Println("Waking up worker " + strconv.Itoa(i))
		pizzaWorkers[i] = components.PizzaWorker{Name: uint64(i + 1), HasAssignedOven: false}
		wg.Add(1)
	}
	log.Printf("Starting the pizzeria")
	start := time.Now()
	for i := 0; i < configs.Parameters.NumberOfWorker; i++ {
		go pizzaWorkers[i].Work(&wg)
	}
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Took %s to cook %d pizzas with %d workers and %d ovens", elapsed, configs.Parameters.NumberOfOrder, configs.Parameters.NumberOfWorker, configs.Parameters.NumberOfOven)
}
