// This file is the entrypoint to the pizzeria
package components

import (
	"fmt"
	config "pizzago/internal/config"
	"strconv"
	"sync"
)

// Type to describe a pizza that can either be baked or not.
type Pizza struct {
	isBaked bool
}

// The config that will be used for the run
var Config config.Config

// The slice containing all the pizza workers
var pizzaWorkers []PizzaWorker

// The slice containing all the ovens
var ovenList []PizzaOven

// Initialize a slice of ovens of size 'NumberOfOvens'
func InitOvens() []PizzaOven {
	var ovenList = make([]PizzaOven, Config.Parameters.NumberOfOvens)
	for i := uint64(0); i < Config.Parameters.NumberOfOvens; i++ {
		ovenList[i] = PizzaOven{isUsed: 0}
	}
	return ovenList
}

// Initialize a slice of bakers of size 'NumberOfWorkers', the uid of the bakers
// goes from 1 to NumberOfWorkers
func InitBakers() []PizzaWorker {
	var pizzaWorkers = make([]PizzaWorker, Config.Parameters.NumberOfWorkers)
	for i := range pizzaWorkers {
		fmt.Println("Waking up worker " + strconv.Itoa(i))
		pizzaWorkers[i] = PizzaWorker{Name: uint64(i + 1), HasAssignedOven: false}
	}
	return pizzaWorkers
}

// The main function of the pizzeria, used to start it.
func StartPizzeria() {
	// Read the configuration specified
	config.ReadConfig(&Config)
	// Initialize the ovens
	ovenList = InitOvens()
	// Initialize the bakers
	pizzaWorkers = InitBakers()

	// A waitGroup that will be helpful to wait for all
	// bakers before returning
	var wg sync.WaitGroup

	// Start a GoRoutine for each baker
	for _, pizzaWorker := range pizzaWorkers {
		wg.Add(1)
		go pizzaWorker.Work(&wg)
	}
	// Wait for all bakers
	wg.Wait()
}
