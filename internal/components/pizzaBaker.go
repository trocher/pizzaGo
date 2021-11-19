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

// Interface of a PizzaBaker
type PizzaBaker interface {
	ProcessOrder() Order
	Prepare(order Order) *Pizza
	QualityCheck(pizza *Pizza) (*Pizza, error)
	Work()
	FindOven() *Oven
	ReleaseOven(o *Oven)
}

type PizzaWorker struct {
	Name            uint64
	HasAssignedOven bool
}

func (w *PizzaWorker) ProcessOrder() Order {
	orderId := atomic.AddUint64(&OrderTaken, 1)
	timing.WaitFor(configs.Timings.Process)
	return Order{orderId}
}

func (w *PizzaWorker) Prepare(order Order) *Pizza {
	timing.WaitFor(configs.Timings.Prepare)
	return &Pizza{false}
}

func (w *PizzaWorker) QualityCheck(pizza *Pizza) (*Pizza, error) {
	timing.WaitFor(configs.Timings.QualityCheck)
	if !pizza.isBaked {
		return pizza, errors.New("Pizza worker " + strconv.FormatUint(w.Name, 10) + " : Sorry I forgot to bake your pizza")
	}
	return pizza, nil

}

func (w *PizzaWorker) FindOven() *PizzaOven {
	if w.HasAssignedOven {
		return &ovenList[w.Name-1]
	}
	for {
		for i := range ovenList {
			if atomic.CompareAndSwapUint64(&(ovenList[i].isUsed), 0, w.Name) {
				return &ovenList[i]
			}
		}
	}
}

func (w *PizzaWorker) ReleaseOven(o *PizzaOven) {
	if w.HasAssignedOven {
		return
	}
	if !atomic.CompareAndSwapUint64(&(o.isUsed), w.Name, 0) {
		log.Printf("Pizza worker " + strconv.FormatUint(w.Name, 10) + " : Someone is using my oven")
	}
}

func (w *PizzaWorker) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	if configs.Parameters.NumberOfOven >= uint64(configs.Parameters.NumberOfWorker) {
		w.HasAssignedOven = true
		ovenList[w.Name-1].isUsed = w.Name
	}
	for atomic.LoadUint64(&OrderTaken) < configs.Parameters.NumberOfOrder {
		start := time.Now()

		order := w.ProcessOrder()
		oven := w.FindOven()
		pizza := w.Prepare(order)
		*pizza = oven.Bake(*pizza)
		w.ReleaseOven(oven)
		pizza, error := w.QualityCheck(pizza)

		if error != nil {
			log.Fatal(error)
		}

		atomic.AddUint64(&OrderDelivered, 1)

		elapsed := time.Since(start)
		log.Printf("Order %d took %s", order.id, elapsed)
	}
}
