package components

type Oven interface {
	Bake(unbakedPizza Pizza) Pizza
}

type PizzaOven struct {
	isUsed uint64
}
