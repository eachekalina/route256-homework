package packaging

import (
	"homework/internal/app/order"
)

type WrapVariant struct {
}

func (v WrapVariant) Apply(o order.Order) (order.Order, error) {
	o.PriceRub += 1
	return o, nil
}
