// This file describes a pizza baker
package components

import (
	"errors"
	"log"
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

// A pizza worker has a name between 1 and NumberOfWorkers
// The flag HasAssignedOven is true if the worker has an assigned
// oven during the shift. He also hold pointers to usefull slices
// and integers
type PizzaWorker struct {
	Name                  uint64
	HasAssignedOven       bool
	orderTaken            *uint64
	ovenList              *[]PizzaOven
	timeTakenTakingOrders *uint64
}

// Process an order
func (w *PizzaWorker) ProcessOrder() (Order, error) {
	// Increment atomically the counter of taken order
	orderId := atomic.AddUint64(w.orderTaken, 1)
	// If the order that has been taken is above the limit, cancel it.
	if orderId > Config.Parameters.NumberOfOrders {
		atomic.SwapUint64(w.orderTaken, Config.Parameters.NumberOfOrders)
		return Order{orderId}, errors.New("Worker " + strconv.FormatUint(w.Name, 10) + " Oups, I took too much order, I won't do this one")
	}
	WaitFor(Config.Times.Process)
	return Order{orderId}, nil
}

// Prepare an order
func (w *PizzaWorker) Prepare(order Order) *Pizza {
	WaitFor(Config.Times.Prepare)
	return &Pizza{false}
}

// Check the quality of a pizza
func (w *PizzaWorker) QualityCheck(pizza *Pizza) (*Pizza, error) {
	WaitFor(Config.Times.QualityCheck)
	if !pizza.isBaked {
		return pizza, errors.New("Worker " + strconv.FormatUint(w.Name, 10) + " : Sorry I forgot to bake your pizza")
	}
	return pizza, nil

}

// Find an oven, if there are more worker than oven, the worker does not have an
// assigned ove, he has to go throught the list of oven to find an unused one.
// All oven managment is done using atomic operations. A worker can claim an oven
// by writting his uid in the isUsed field.
func (w *PizzaWorker) FindOven() *PizzaOven {
	if w.HasAssignedOven {
		return &(*w.ovenList)[w.Name-1]
	}
	for {
		for i := range *(w.ovenList) {
			if atomic.CompareAndSwapUint64(&((*(w.ovenList))[i].isUsed), 0, w.Name) {
				return &(*w.ovenList)[i]
			}
		}
		WaitFor(1)
	}
}

// Release an oven, simply write back 0 in the isUsed field of the oven to notify
// other worker that the oven is no longer in use
func (w *PizzaWorker) ReleaseOven(o *PizzaOven) {
	if w.HasAssignedOven {
		return
	}
	if !atomic.CompareAndSwapUint64(&(o.isUsed), w.Name, 0) {
		log.Printf("Worker " + strconv.FormatUint(w.Name, 10) + " : Someone is using my oven")
	}
}

// Main function of a worker
func (w PizzaWorker) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	// Loop that goes on while there are still orders to be taken care of
	for atomic.LoadUint64(w.orderTaken) < Config.Parameters.NumberOfOrders {
		order, error := w.ProcessOrder()
		// If the worker managed to enter this loop with other workers when there was
		// not enought orders lefts for all of them, he cancel his order and break the loop
		if error != nil {
			break
		}

		pizza := w.Prepare(order)
		oven := w.FindOven()
		*pizza = oven.Bake(*pizza)
		w.ReleaseOven(oven)

		pizza, error = w.QualityCheck(pizza)

		// If the quality check is not successful, there must be an error
		// in the program somewhere
		if error != nil {
			log.Fatal(error)
		}

	}
	elapsed := time.Since(start)
	atomic.AddUint64(w.timeTakenTakingOrders, uint64(elapsed))
}
