package components

import (
	"log"
	"pizzago/internal/configs"
	"pizzago/internal/timing"
)

type Oven interface {
	Bake(unbakedPizza Pizza) Pizza
}

type PizzaOven struct {
	isUsed uint64
}

func (o *PizzaOven) Bake(unbakedPizza Pizza) Pizza {
	if unbakedPizza.isBaked {
		log.Printf("That is weird this pizza is already baked")
		return unbakedPizza
	}
	timing.WaitFor(configs.Timings.Bake)
	unbakedPizza.isBaked = true
	return unbakedPizza
}
