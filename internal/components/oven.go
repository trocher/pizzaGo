// This package provide the definition and implementation of an oven
package components

import (
	"log"
)

// Iterface of an oven
type Oven interface {
	Bake(unbakedPizza Pizza) Pizza
}

// Implementation of an oven, isUsed is used as a lock by workers, i.e, when it is
// is equal to 0, it is free to use, otherwise, a worker will set th value to its
// own uid to specify that he is using it
type PizzaOven struct {
	isUsed uint64
}

// The method bake
func (o *PizzaOven) Bake(unbakedPizza Pizza) Pizza {
	// Check wether the pizza is already backed, which should not happen
	if unbakedPizza.isBaked {
		log.Fatal("That is weird this pizza is already baked")
	}
	// Wait for the time needed to bake a pizza
	WaitFor(Config.Times.Bake)
	// Set the isBaked flag to true
	unbakedPizza.isBaked = true
	return unbakedPizza
}
