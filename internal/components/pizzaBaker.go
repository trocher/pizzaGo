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
// oven during the shift.
type PizzaWorker struct {
	Name            uint64
	HasAssignedOven bool
}

// Process an order
func (w *PizzaWorker) ProcessOrder() (Order, error) {
	// Increment atomically the counter of taken order
	orderId := atomic.AddUint64(&OrderTaken, 1)
	// If the order that has been taken is above the limit, cancel it.
	if orderId > Config.Parameters.NumberOfOrders {
		atomic.SwapUint64(&OrderTaken, Config.Parameters.NumberOfOrders)
		return Order{orderId}, errors.New("Oups, I took too much order")
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
		return pizza, errors.New("Pizza worker " + strconv.FormatUint(w.Name, 10) + " : Sorry I forgot to bake your pizza")
	}
	return pizza, nil

}

// Find an oven, if there are more worker than oven, the worker does not have an
// assigned ove, he has to go throught the list of oven to find an unused one.
// All oven managment is done using atomic operations. A worker can claim an oven
// by writting his uid in the isUsed field.
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

// Release an oven, simply write back 0 in the isUsed field of the oven to notify
// other worker that the oven is no longer in use
func (w *PizzaWorker) ReleaseOven(o *PizzaOven) {
	if w.HasAssignedOven {
		return
	}
	if !atomic.CompareAndSwapUint64(&(o.isUsed), w.Name, 0) {
		log.Printf("Pizza worker " + strconv.FormatUint(w.Name, 10) + " : Someone is using my oven")
	}
}

// Main function of a worker
func (w *PizzaWorker) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	// If there are less workers than ovens, then the worker can claim an oven, he
	// clain the oven w.Name-1
	if Config.Parameters.NumberOfOvens >= uint64(Config.Parameters.NumberOfWorkers) {
		w.HasAssignedOven = true
		ovenList[w.Name-1].isUsed = w.Name
	}

	// Loop that goes on while there are still orders to be taken care of
	for atomic.LoadUint64(&OrderTaken) < Config.Parameters.NumberOfOrders {
		start := time.Now()

		order, error := w.ProcessOrder()
		// If the worker managed to enter this loop with other workers when there was
		// not enought orders lefts for all of them, he cancel his order and break the loop
		if error != nil {
			log.Println(error)
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

		// Increment the number of delivered order, used to check the correcness
		// of the program
		atomic.AddUint64(&OrderDelivered, 1)

		elapsed := time.Since(start)
		log.Printf("Order %d took %s", order.id, elapsed)
	}
}
