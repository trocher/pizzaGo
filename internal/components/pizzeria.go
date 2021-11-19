package components

import (
	"fmt"
	config "pizzago/internal/configs"
	"strconv"
	"sync"
)

// === Pizza definition ===
type Pizza struct {
	isBaked bool
}

var Config config.Config
var pizzaWorkers []PizzaWorker
var ovenList []PizzaOven

func InitOvens() []PizzaOven {
	var ovenList = make([]PizzaOven, Config.Parameters.NumberOfOvens)
	for i := uint64(0); i < Config.Parameters.NumberOfOvens; i++ {
		ovenList[i] = PizzaOven{isUsed: 0}
	}
	return ovenList
}

func InitBakers() []PizzaWorker {
	var pizzaWorkers = make([]PizzaWorker, Config.Parameters.NumberOfWorkers)
	for i := range pizzaWorkers {
		fmt.Println("Waking up worker " + strconv.Itoa(i))
		pizzaWorkers[i] = PizzaWorker{Name: uint64(i + 1), HasAssignedOven: false}
	}
	return pizzaWorkers
}

func StartPizzeria() {
	config.ReadConfig(&Config)
	ovenList = InitOvens()
	pizzaWorkers = InitBakers()
	fmt.Println("Waking up waozd")

	var wg sync.WaitGroup

	for _, pizzaWorker := range pizzaWorkers {
		wg.Add(1)
		go pizzaWorker.Work(&wg)
	}
	wg.Wait()

}
