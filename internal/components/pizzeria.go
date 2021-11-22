// This file is the entrypoint to the pizzeria
package components

import (
	"log"
	"pizzago/internal/config"
	"sync"
	"time"
)

// Type to describe a pizza that can either be baked or not.
type Pizza struct {
	isBaked bool
}

// An order is defined by its id
type Order struct {
	id uint64
}

// The config that will be used for the run
var Config config.Config

// Initialize a slice of ovens of size 'NumberOfOvens'
func InitOvens() []PizzaOven {
	ovenList := make([]PizzaOven, Config.Parameters.NumberOfOvens)
	for i := uint64(0); i < Config.Parameters.NumberOfOvens; i++ {
		ovenList[i] = PizzaOven{isUsed: 0}
	}
	return ovenList
}

// Initialize a slice of bakers of size 'NumberOfWorkers', the uid of the bakers
// goes from 1 to NumberOfWorkers, they also have pointers to the list of oven,
// and the orderTaken counter
func InitBakers(hasAssignedOven bool, ovenList *[]PizzaOven, orderTaken *uint64, timeTakenTakingOrders *uint64) []PizzaWorker {
	pizzaWorkers := make([]PizzaWorker, Config.Parameters.NumberOfWorkers)
	for i := range pizzaWorkers {
		pizzaWorkers[i] = PizzaWorker{Name: uint64(i + 1), HasAssignedOven: hasAssignedOven, ovenList: ovenList, orderTaken: orderTaken, timeTakenTakingOrders: timeTakenTakingOrders}
	}
	return pizzaWorkers
}

// The main function of the pizzeria, used to start it.
func StartPizzeria(withConfig config.Config) time.Duration {
	// Read the configuration specified
	Config = withConfig

	// ==== Initializations ====

	ovenList := InitOvens()

	var orderTaken uint64 = 0
	var timeTakenTakingOrders uint64 = 0

	// If there are less workers than ovens, then the worker can claim an oven, he
	// claim the oven w.Name-1
	hasAssignedOven := false
	if Config.Parameters.NumberOfOvens >= uint64(Config.Parameters.NumberOfWorkers) {
		hasAssignedOven = true
	}

	pizzaWorkers := InitBakers(hasAssignedOven, &ovenList, &orderTaken, &timeTakenTakingOrders)

	// ==== Start the bakers =====

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

	if orderTaken != Config.Parameters.NumberOfOrders {
		log.Fatal("The number of order taken is different from the number of delivered orders")
	}

	// compute and return the average latency
	return time.Duration(float64(timeTakenTakingOrders) / float64(Config.Parameters.NumberOfOrders))
}
