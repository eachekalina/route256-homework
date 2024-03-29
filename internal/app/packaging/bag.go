package packaging

import (
	"errors"
	"homework/internal/app/order"
)

type BagVariant struct {
}

func (v BagVariant) Apply(o order.Order) (order.Order, error) {
	if o.WeightKg >= 10.0 {
		return order.Order{}, errors.New("bag cannot handle more than 10 kg")
	}

	o.PriceRub += 5
	return o, nil
}
