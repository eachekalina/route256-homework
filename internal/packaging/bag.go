package packaging

import (
	"errors"
	"homework/internal/model"
)

type BagVariant struct {
}

func (v BagVariant) Apply(order model.Order) (model.Order, error) {
	if order.WeightKg >= 10.0 {
		return model.Order{}, errors.New("bag cannot handle more than 10 kg")
	}

	order.PriceRub += 5
	return order, nil
}
