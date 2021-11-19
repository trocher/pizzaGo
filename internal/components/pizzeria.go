package components

import (
	"pizzago/internal/configs"
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
