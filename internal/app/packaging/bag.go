package packaging

import (
	"errors"
	"homework/internal/app/order"
)

type Bag struct {
}

func (v Bag) Apply(o order.Order) (order.Order, error) {
	if o.WeightKg >= 10.0 {
		return o, errors.New("bag cannot handle more than 10 kg")
	}

	o.PriceRub += 5
	return o, nil
}
