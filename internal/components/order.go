// This file define the implementation of an order

package components

// An order is defined by its id
type Order struct {
	id uint64
}

// A counter that is incremented atomically each time a worker takes a new order
// Represent hence the flow of clients
var OrderTaken uint64 = 0

// An atomic counter that is incremented atomically each time an order is finished
// and the pizza has been given to the customer
var OrderDelivered uint64 = 0
