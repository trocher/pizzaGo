package components

import (
	"errors"
	"log"
	"pizzago/internal/configs"
	"pizzago/internal/timing"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type PizzaBaker interface {
	ProcessOrder() Order
	Prepare(order Order, oven *PizzaOven) *Pizza
	QualityCheck(pizza *Pizza) (*Pizza, error)
	Work()
}

type PizzaWorker struct {
	Name uint64
}

func (w *PizzaWorker) ProcessOrder() Order {
	orderId := atomic.AddUint64(&orderTaken, 1)
	timing.WaitFor(configs.Timings.Process)
	return Order{orderId}
}

func (w *PizzaWorker) Prepare(order Order, oven *PizzaOven) *Pizza {
	timing.WaitFor(configs.Timings.Prepare)
	timing.WaitFor(configs.Timings.Bake)
	return &Pizza{true}
}

func (w *PizzaWorker) QualityCheck(pizza *Pizza) (*Pizza, error) {
	timing.WaitFor(configs.Timings.QualityCheck)
	if !pizza.isBaked {
		return pizza, errors.New("Pizza worker " + strconv.FormatUint(w.Name, 10) + " : Sorry I forgot to bake your pizza")
	}
	return pizza, nil

}

func (w *PizzaWorker) FindOven() *PizzaOven {
	for {
		for i := range ovenList {
			if atomic.CompareAndSwapUint64(&(ovenList[i].isUsed), 0, w.Name) {
				return &ovenList[i]
			}
		}
	}
}

func (w *PizzaWorker) ReleaseOven(o *PizzaOven) {
	if !atomic.CompareAndSwapUint64(&(o.isUsed), w.Name, 0) {
		log.Printf("Pizza worker " + strconv.FormatUint(w.Name, 10) + " : Someone is using my oven")
	}
}

func (w *PizzaWorker) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	for atomic.LoadUint64(&orderTaken) < configs.Parameters.NumberOfOrder {
		start := time.Now()

		order := w.ProcessOrder()
		oven := w.FindOven()
		pizza := w.Prepare(order, oven)
		w.ReleaseOven(oven)
		pizza, error := w.QualityCheck(pizza)

		if error != nil {
			log.Fatal(error)
		}

		atomic.AddUint64(&orderDelivered, 1)

		elapsed := time.Since(start)
		log.Printf("Order %d took %s", order.id, elapsed)
	}
}
