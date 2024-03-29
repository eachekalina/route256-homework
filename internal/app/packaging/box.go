package packaging

import (
	"errors"
	"homework/internal/app/order"
)

type BoxVariant struct {
}

func (v BoxVariant) Apply(o order.Order) (order.Order, error) {
	if o.WeightKg >= 30.0 {
		return order.Order{}, errors.New("box cannot handle more than 30 kg")
	}

	o.PriceRub += 20
	return o, nil
}
