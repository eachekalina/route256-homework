package packaging

import (
	"homework/internal/model"
)

type WrapVariant struct {
}

func (v WrapVariant) Apply(order model.Order) (model.Order, error) {
	order.PriceRub += 1
	return order, nil
}
