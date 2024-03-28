package packaging

import (
	"errors"
	"homework/internal/model"
)

type BoxVariant struct {
}

func (v BoxVariant) Apply(order model.Order) (model.Order, error) {
	if order.WeightKg >= 30.0 {
		return model.Order{}, errors.New("box cannot handle more than 30 kg")
	}

	order.PriceRub += 20
	return order, nil
}
