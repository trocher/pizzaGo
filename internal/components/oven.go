// This file provide the definition and implementation of an oven
package components

// Interface of an oven
type Oven interface {
	Bake(unbakedPizza Pizza) Pizza
}

// Implementation of an oven, isUsed is used as a lock by workers, i.e, when it is
// equal to 0, it is free to use, otherwise, a worker will set the value to its
// own uid to indicate that he is using it.
type PizzaOven struct {
	isUsed uint64
}

// The method bake
func (o *PizzaOven) Bake(unbakedPizza Pizza) Pizza {
	// Wait for the time needed to bake a pizza
	WaitFor(Config.Times.Bake)
	// Set the isBaked flag to true
	unbakedPizza.isBaked = true
	return unbakedPizza
}
