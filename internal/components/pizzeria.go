package components

import (
	"fmt"
	"pizzago/internal/configs"
	"pizzago/internal/timing"
)

// === Pizza definition ===
type Pizza struct {
	isBaked bool
}

var ovenList = make([]PizzaOven, configs.Parameters.NumberOfOven)

func InitOven() {
	for i := range ovenList {
		ovenList[i] = PizzaOven{isUsed: 0}
	}
}

func (o PizzaOven) Bake(unbakedPizza Pizza) Pizza {
	if unbakedPizza.isBaked {
		fmt.Println("That is weird, this pizza is already baked")
		return unbakedPizza
	}
	timing.WaitFor(configs.Timings.Bake)
	unbakedPizza.isBaked = true
	return unbakedPizza
}

// === Pizza backer ===
