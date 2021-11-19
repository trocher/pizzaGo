package components

import (
	"fmt"
	"log"
	"pizzago/internal/configs"
	"strconv"
	"sync"
)

// === Pizza definition ===
type Pizza struct {
	isBaked bool
}

var ovenList = make([]PizzaOven, configs.Parameters.NumberOfOven)
var pizzaWorkers = make([]PizzaWorker, configs.Parameters.NumberOfWorker)

func InitOven() {
	for i := range ovenList {
		ovenList[i] = PizzaOven{isUsed: 0}
	}
}

func InitBakers() {
	for i := range pizzaWorkers {
		fmt.Println("Waking up worker " + strconv.Itoa(i))
		pizzaWorkers[i] = PizzaWorker{Name: uint64(i + 1), HasAssignedOven: false}
	}
}

func StartPizzeria() {
	InitOven()
	InitBakers()
	var wg sync.WaitGroup
	log.Printf("Starting the pizzeria")

	for _, pizzaWorker := range pizzaWorkers {
		wg.Add(1)
		go pizzaWorker.Work(&wg)
	}
	wg.Wait()

}
