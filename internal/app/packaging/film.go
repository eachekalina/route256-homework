package packaging

import (
	"homework/internal/app/order"
)

type Film struct {
}

func (v Film) Apply(o order.Order) (order.Order, error) {
	o.PriceRub += 1
	return o, nil
}
