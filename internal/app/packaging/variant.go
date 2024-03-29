package packaging

import (
	"homework/internal/app/order"
)

type Type string

const (
	BagType  Type = "bag"
	BoxType  Type = "box"
	WrapType Type = "wrap"
)

type Variant interface {
	Apply(order order.Order) (order.Order, error)
}
